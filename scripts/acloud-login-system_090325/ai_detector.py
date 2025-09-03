"""
AI Element Detector Module
Uses GPT-4 Vision to detect form elements on web pages
"""

import asyncio
import base64
import json
from typing import Dict, Any, Optional, List
import openai
from playwright.async_api import Page

from config.settings import Settings
from modules.logger import setup_logger


class AIElementDetector:
    """AI-powered element detection using GPT-4 Vision"""
    
    def __init__(self, config: Settings):
        self.config = config
        self.logger = setup_logger()
        self.client = openai.AsyncOpenAI(api_key=config.OPENAI_API_KEY)
        
    async def detect_elements(self, page: Page) -> Dict[str, Any]:
        """
        Detect form elements using GPT-4 Vision
        
        Args:
            page: Playwright page object
            
        Returns:
            Dict containing detected elements
        """
        try:
            # Take screenshot
            screenshot = await page.screenshot(full_page=False, type='png')
            screenshot_b64 = base64.b64encode(screenshot).decode('utf-8')
            
            # Analyze with GPT-4 Vision
            element_coords = await self._analyze_screenshot(screenshot_b64)
            
            if not element_coords:
                return {}
            
            # Convert coordinates to selectors
            elements = await self._coords_to_elements(page, element_coords)
            
            return elements
            
        except Exception as e:
            self.logger.error(f"AI element detection failed: {str(e)}")
            return {}
    
    async def _analyze_screenshot(self, screenshot_b64: str) -> Optional[Dict[str, Dict[str, int]]]:
        """
        Send screenshot to GPT-4 Vision for analysis
        
        Args:
            screenshot_b64: Base64 encoded screenshot
            
        Returns:
            Dictionary with element coordinates
        """
        try:
            prompt = """
            Analyze this login page screenshot and identify the following form elements:
            1. Email/username input field
            2. Password input field  
            3. Submit/login button
            4. Any CAPTCHA elements
            
            For each element found, provide the approximate center coordinates (x, y) as percentages of the image dimensions (0-100).
            
            Return your response as a JSON object with this exact structure:
            {
                "email_field": {"x": 50, "y": 30},
                "password_field": {"x": 50, "y": 45},
                "submit_button": {"x": 50, "y": 60},
                "captcha_element": {"x": 50, "y": 75}
            }
            
            Only include elements that are clearly visible. If an element is not found, omit it from the JSON.
            Ensure coordinates are realistic for a typical login form layout.
            """
            
            response = await self.client.chat.completions.create(
                model=self.config.OPENAI_MODEL,
                messages=[{
                    "role": "user",
                    "content": [
                        {"type": "text", "text": prompt},
                        {
                            "type": "image_url",
                            "image_url": {
                                "url": f"data:image/png;base64,{screenshot_b64}",
                                "detail": "high"
                            }
                        }
                    ]
                }],
                max_tokens=self.config.OPENAI_MAX_TOKENS,
                temperature=0.1
            )
            
            content = response.choices[0].message.content
            self.logger.debug(f"GPT-4V response: {content}")
            
            # Extract JSON from response
            json_start = content.find('{')
            json_end = content.rfind('}') + 1
            
            if json_start != -1 and json_end > json_start:
                json_str = content[json_start:json_end]
                coords = json.loads(json_str)
                self.logger.info(f"Detected element coordinates: {coords}")
                return coords
            else:
                self.logger.warning("No valid JSON found in GPT-4V response")
                return None
                
        except json.JSONDecodeError as e:
            self.logger.error(f"Failed to parse GPT-4V JSON response: {str(e)}")
            return None
        except Exception as e:
            self.logger.error(f"GPT-4V analysis failed: {str(e)}")
            return None
    
    async def _coords_to_elements(self, page: Page, coords: Dict[str, Dict[str, int]]) -> Dict[str, Any]:
        """
        Convert percentage coordinates to actual page elements
        
        Args:
            page: Playwright page object
            coords: Element coordinates as percentages
            
        Returns:
            Dictionary of detected elements
        """
        elements = {}
        
        # Get viewport size
        viewport = page.viewport_size
        if not viewport:
            viewport = {'width': self.config.VIEWPORT_WIDTH, 'height': self.config.VIEWPORT_HEIGHT}
        
        width = viewport['width']
        height = viewport['height']
        
        # Convert percentage coordinates to pixel coordinates
        for element_type, coord in coords.items():
            try:
                x = int((coord['x'] / 100) * width)
                y = int((coord['y'] / 100) * height)
                
                self.logger.debug(f"Looking for {element_type} at coordinates ({x}, {y})")
                
                # Find element at coordinates
                element = await self._find_element_at_coordinates(page, x, y, element_type)
                
                if element:
                    # Map to standard names
                    standard_name = self._map_element_name(element_type)
                    elements[standard_name] = element
                    self.logger.info(f"Found {standard_name} element")
                else:
                    self.logger.warning(f"No element found at {element_type} coordinates")
                    
            except Exception as e:
                self.logger.error(f"Error processing {element_type} coordinates: {str(e)}")
        
        return elements
    
    async def _find_element_at_coordinates(self, page: Page, x: int, y: int, element_type: str) -> Optional[Any]:
        """
        Find the best matching element near the given coordinates
        
        Args:
            page: Playwright page object
            x, y: Pixel coordinates
            element_type: Type of element being searched
            
        Returns:
            Playwright element or None
        """
        try:
            # First try exact coordinates
            element = await page.evaluate(f"""
                () => {{
                    const element = document.elementFromPoint({x}, {y});
                    return element ? element.tagName + '_' + (element.type || 'unknown') : null;
                }}
            """)
            
            if element:
                # Get the actual element
                locator = page.locator(f'xpath=//*[@data-playwright-target]').first
                
                # Try different strategies based on element type
                if 'email' in element_type.lower():
                    candidates = await self._find_email_candidates(page, x, y)
                elif 'password' in element_type.lower():
                    candidates = await self._find_password_candidates(page, x, y)
                elif 'submit' in element_type.lower() or 'button' in element_type.lower():
                    candidates = await self._find_button_candidates(page, x, y)
                elif 'captcha' in element_type.lower():
                    candidates = await self._find_captcha_candidates(page, x, y)
                else:
                    candidates = []
                
                if candidates:
                    # Return the closest candidate
                    return candidates[0]
            
            return None
            
        except Exception as e:
            self.logger.error(f"Error finding element at coordinates: {str(e)}")
            return None
    
    async def _find_email_candidates(self, page: Page, x: int, y: int) -> List[Any]:
        """Find email input field candidates near coordinates"""
        selectors = [
            'input[type="email"]',
            'input[type="text"][name*="email" i]',
            'input[type="text"][placeholder*="email" i]',
            'input[type="text"][id*="email" i]',
            'input[autocomplete="email"]',
            'input[autocomplete="username"]'
        ]
        
        candidates = []
        for selector in selectors:
            try:
                elements = page.locator(selector)
                count = await elements.count()
                
                for i in range(count):
                    element = elements.nth(i)
                    if await element.is_visible():
                        # Check if element is near target coordinates
                        box = await element.bounding_box()
                        if box and self._is_near_coordinates(box, x, y, tolerance=100):
                            candidates.append(element)
            except:
                continue
        
        return candidates
    
    async def _find_password_candidates(self, page: Page, x: int, y: int) -> List[Any]:
        """Find password input field candidates near coordinates"""
        selectors = [
            'input[type="password"]',
            'input[autocomplete="current-password"]',
            'input[autocomplete="password"]'
        ]
        
        candidates = []
        for selector in selectors:
            try:
                elements = page.locator(selector)
                count = await elements.count()
                
                for i in range(count):
                    element = elements.nth(i)
                    if await element.is_visible():
                        box = await element.bounding_box()
                        if box and self._is_near_coordinates(box, x, y, tolerance=100):
                            candidates.append(element)
            except:
                continue
        
        return candidates
    
    async def _find_button_candidates(self, page: Page, x: int, y: int) -> List[Any]:
        """Find submit button candidates near coordinates"""
        selectors = [
            'button[type="submit"]',
            'input[type="submit"]',
            'button:has-text("Sign in")',
            'button:has-text("Login")',
            'button:has-text("Submit")',
            '[data-testid*="submit"]',
            '[data-testid*="login"]',
            '.submit-button',
            '.login-button'
        ]
        
        candidates = []
        for selector in selectors:
            try:
                elements = page.locator(selector)
                count = await elements.count()
                
                for i in range(count):
                    element = elements.nth(i)
                    if await element.is_visible():
                        box = await element.bounding_box()
                        if box and self._is_near_coordinates(box, x, y, tolerance=150):
                            candidates.append(element)
            except:
                continue
        
        return candidates
    
    async def _find_captcha_candidates(self, page: Page, x: int, y: int) -> List[Any]:
        """Find CAPTCHA element candidates near coordinates"""
        selectors = [
            '.g-recaptcha',
            '[data-sitekey]',
            '.recaptcha',
            '.captcha',
            '.h-captcha',
            '[data-hcaptcha-sitekey]',
            'iframe[src*="recaptcha"]',
            'iframe[src*="hcaptcha"]'
        ]
        
        candidates = []
        for selector in selectors:
            try:
                elements = page.locator(selector)
                count = await elements.count()
                
                for i in range(count):
                    element = elements.nth(i)
                    if await element.is_visible():
                        box = await element.bounding_box()
                        if box and self._is_near_coordinates(box, x, y, tolerance=200):
                            candidates.append(element)
            except:
                continue
        
        return candidates
    
    def _is_near_coordinates(self, box: Dict[str, float], target_x: int, target_y: int, tolerance: int = 100) -> bool:
        """Check if element bounding box is near target coordinates"""
        element_center_x = box['x'] + box['width'] / 2
        element_center_y = box['y'] + box['height'] / 2
        
        distance = ((element_center_x - target_x) ** 2 + (element_center_y - target_y) ** 2) ** 0.5
        
        return distance <= tolerance
    
    def _map_element_name(self, element_type: str) -> str:
        """Map AI-detected element types to standard names"""
        mapping = {
            'email_field': 'email',
            'username_field': 'email',
            'password_field': 'password',
            'submit_button': 'submit',
            'login_button': 'submit',
            'captcha_element': 'captcha',
            'recaptcha_element': 'captcha'
        }
        
        return mapping.get(element_type, element_type)