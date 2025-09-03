"""
Element Mapper Module
Fallback element detection using traditional methods
"""

import asyncio
from typing import Dict, Any, List, Optional
from playwright.async_api import Page

from config.settings import Settings
from modules.logger import setup_logger


class ElementMapper:
    """Traditional element detection with multiple fallback strategies"""
    
    def __init__(self, config: Settings):
        self.config = config
        self.logger = setup_logger()
    
    async def detect_elements_fallback(self, page: Page) -> Dict[str, Any]:
        """
        Detect form elements using traditional methods as fallback
        
        Args:
            page: Playwright page object
            
        Returns:
            Dict containing detected elements
        """
        self.logger.info("Starting fallback element detection")
        
        # Try multiple detection strategies
        strategies = [
            self._semantic_detection,
            self._heuristic_detection,
            self._aria_detection,
            self._fuzzy_detection,
            self._last_resort_detection
        ]
        
        for i, strategy in enumerate(strategies):
            try:
                self.logger.debug(f"Trying detection strategy {i+1}")
                elements = await strategy(page)
                
                if self._validate_elements(elements):
                    self.logger.info(f"Strategy {i+1} successful")
                    return elements
                else:
                    self.logger.debug(f"Strategy {i+1} found incomplete elements")
                    
            except Exception as e:
                self.logger.warning(f"Strategy {i+1} failed: {str(e)}")
                continue
        
        self.logger.warning("All detection strategies failed")
        return {}
    
    async def _semantic_detection(self, page: Page) -> Dict[str, Any]:
        """Modern semantic HTML detection"""
        elements = {}
        
        # Email/username field detection
        email_selectors = [
            'input[type="email"]',
            'input[autocomplete="email"]',
            'input[autocomplete="username"]',
            'input[name="email" i]',
            'input[name="username" i]',
            'input[name="user" i]',
            'input[id*="email" i]',
            'input[id*="username" i]',
            'input[placeholder*="email" i]',
            'input[placeholder*="username" i]'
        ]
        
        for selector in email_selectors:
            try:
                element = page.locator(selector).first
                if await element.is_visible():
                    elements['email'] = element
                    self.logger.debug(f"Found email field with selector: {selector}")
                    break
            except:
                continue
        
        # Password field detection
        password_selectors = [
            'input[type="password"]',
            'input[autocomplete="current-password"]',
            'input[autocomplete="password"]',
            'input[name="password" i]',
            'input[id*="password" i]',
            'input[placeholder*="password" i]'
        ]
        
        for selector in password_selectors:
            try:
                element = page.locator(selector).first
                if await element.is_visible():
                    elements['password'] = element
                    self.logger.debug(f"Found password field with selector: {selector}")
                    break
            except:
                continue
        
        # Submit button detection
        submit_selectors = [
            'button[type="submit"]',
            'input[type="submit"]',
            'button:has-text("Sign in") >> visible=true',
            'button:has-text("Login") >> visible=true',
            'button:has-text("Submit") >> visible=true',
            'button:has-text("Continue") >> visible=true',
            '[data-testid*="submit" i]',
            '[data-testid*="login" i]',
            '[data-testid*="signin" i]',
            'button[class*="submit" i]',
            'button[class*="login" i]',
            'input[value*="sign" i]',
            'input[value*="login" i]'
        ]
        
        for selector in submit_selectors:
            try:
                element = page.locator(selector).first
                if await element.is_visible():
                    elements['submit'] = element
                    self.logger.debug(f"Found submit button with selector: {selector}")
                    break
            except:
                continue
        
        return elements
    
    async def _heuristic_detection(self, page: Page) -> Dict[str, Any]:
        """Pattern-based heuristic detection"""
        elements = {}
        
        # Find forms and analyze their structure
        forms = page.locator('form')
        form_count = await forms.count()
        
        self.logger.debug(f"Found {form_count} forms on page")
        
        for i in range(form_count):
            form = forms.nth(i)
            
            # Skip hidden forms
            if not await form.is_visible():
                continue
            
            # Analyze inputs in this form
            inputs = form.locator('input')
            input_count = await inputs.count()
            
            # Login forms typically have 2-4 inputs
            if 2 <= input_count <= 4:
                form_elements = await self._analyze_form_inputs(form)
                if form_elements:
                    elements.update(form_elements)
                    self.logger.debug(f"Found login form with {input_count} inputs")
                    break
        
        # If no forms found, try global input detection
        if not elements:
            elements = await self._analyze_global_inputs(page)
        
        return elements
    
    async def _analyze_form_inputs(self, form) -> Dict[str, Any]:
        """Analyze inputs within a specific form"""
        elements = {}
        inputs = form.locator('input')
        input_count = await inputs.count()
        
        # Collect input information
        input_info = []
        for i in range(input_count):
            input_elem = inputs.nth(i)
            try:
                if await input_elem.is_visible():
                    input_type = await input_elem.get_attribute('type') or 'text'
                    input_name = await input_elem.get_attribute('name') or ''
                    input_id = await input_elem.get_attribute('id') or ''
                    input_placeholder = await input_elem.get_attribute('placeholder') or ''
                    input_class = await input_elem.get_attribute('class') or ''
                    
                    input_info.append({
                        'element': input_elem,
                        'type': input_type.lower(),
                        'name': input_name.lower(),
                        'id': input_id.lower(),
                        'placeholder': input_placeholder.lower(),
                        'class': input_class.lower()
                    })
            except:
                continue
        
        # Classify inputs
        for info in input_info:
            # Password field
            if info['type'] == 'password':
                elements['password'] = info['element']
            
            # Email field
            elif (info['type'] == 'email' or 
                  'email' in info['name'] or 
                  'email' in info['id'] or 
                  'email' in info['placeholder'] or
                  'username' in info['name'] or
                  'username' in info['id'] or
                  'user' in info['name']):
                elements['email'] = info['element']
            
            # Text field that might be username (if no email field found)
            elif (info['type'] == 'text' and 
                  'email' not in elements and
                  ('user' in info['name'] or 
                   'login' in info['name'] or
                   'account' in info['name'])):
                elements['email'] = info['element']
        
        # Find submit button in form
        submit_buttons = form.locator('button, input[type="submit"]')
        submit_count = await submit_buttons.count()
        
        for i in range(submit_count):
            button = submit_buttons.nth(i)
            try:
                if await button.is_visible():
                    button_text = await button.text_content() or ''
                    button_value = await button.get_attribute('value') or ''
                    button_type = await button.get_attribute('type') or ''
                    
                    if (button_type == 'submit' or
                        'sign' in button_text.lower() or
                        'login' in button_text.lower() or
                        'submit' in button_text.lower() or
                        'continue' in button_text.lower() or
                        'sign' in button_value.lower() or
                        'login' in button_value.lower()):
                        
                        elements['submit'] = button
                        break
            except:
                continue
        
        return elements
    
    async def _analyze_global_inputs(self, page: Page) -> Dict[str, Any]:
        """Analyze all inputs on the page when no suitable form is found"""
        elements = {}
        
        # Find all visible inputs
        all_inputs = page.locator('input:visible')
        input_count = await all_inputs.count()
        
        password_inputs = []
        text_inputs = []
        
        for i in range(input_count):
            input_elem = all_inputs.nth(i)
            try:
                input_type = await input_elem.get_attribute('type') or 'text'
                
                if input_type.lower() == 'password':
                    password_inputs.append(input_elem)
                elif input_type.lower() in ['text', 'email']:
                    text_inputs.append(input_elem)
            except:
                continue
        
        # Assign password field (usually only one)
        if password_inputs:
            elements['password'] = password_inputs[0]
        
        # Assign email field (find the most likely one)
        if text_inputs:
            # Prefer email type or fields with email-related attributes
            email_field = None
            
            for text_input in text_inputs:
                try:
                    input_type = await text_input.get_attribute('type') or 'text'
                    input_name = await text_input.get_attribute('name') or ''
                    input_id = await text_input.get_attribute('id') or ''
                    input_placeholder = await text_input.get_attribute('placeholder') or ''
                    
                    if (input_type.lower() == 'email' or
                        'email' in input_name.lower() or
                        'email' in input_id.lower() or
                        'email' in input_placeholder.lower()):
                        email_field = text_input
                        break
                    elif ('user' in input_name.lower() or
                          'user' in input_id.lower() or
                          'login' in input_name.lower()):
                        email_field = text_input
                except:
                    continue
            
            # If no specific email field found, use first text input
            if not email_field and text_inputs:
                email_field = text_inputs[0]
            
            if email_field:
                elements['email'] = email_field
        
        # Find submit button
        submit_buttons = page.locator('button:visible, input[type="submit"]:visible')
        submit_count = await submit_buttons.count()
        
        for i in range(submit_count):
            button = submit_buttons.nth(i)
            try:
                button_text = await button.text_content() or ''
                button_value = await button.get_attribute('value') or ''
                button_type = await button.get_attribute('type') or ''
                
                if (button_type == 'submit' or
                    'sign' in button_text.lower() or
                    'login' in button_text.lower() or
                    'submit' in button_text.lower() or
                    'continue' in button_text.lower()):
                    
                    elements['submit'] = button
                    break
            except:
                continue
        
        return elements
    
    async def _aria_detection(self, page: Page) -> Dict[str, Any]:
        """ARIA and accessibility-based detection"""
        elements = {}
        
        # ARIA label detection
        aria_selectors = {
            'email': [
                'input[aria-label*="email" i]',
                'input[aria-label*="username" i]',
                'input[aria-label*="user" i]',
                'input[aria-describedby*="email" i]',
                '[role="textbox"][aria-label*="email" i]'
            ],
            'password': [
                'input[aria-label*="password" i]',
                'input[aria-describedby*="password" i]',
                '[role="textbox"][aria-label*="password" i]'
            ],
            'submit': [
                'button[aria-label*="sign" i]',
                'button[aria-label*="login" i]',
                'button[aria-label*="submit" i]',
                '[role="button"][aria-label*="sign" i]',
                '[role="button"][aria-label*="login" i]'
            ]
        }
        
        for element_type, selectors in aria_selectors.items():
            for selector in selectors:
                try:
                    element = page.locator(selector).first
                    if await element.is_visible():
                        elements[element_type] = element
                        self.logger.debug(f"Found {element_type} via ARIA: {selector}")
                        break
                except:
                    continue
        
        return elements
    
    async def _fuzzy_detection(self, page: Page) -> Dict[str, Any]:
        """Fuzzy matching based on common patterns"""
        elements = {}
        
        # Fuzzy text matching for labels and surrounding text
        try:
            # Find inputs near "email" text
            email_candidates = await page.evaluate("""
                () => {
                    const candidates = [];
                    const emailKeywords = ['email', 'e-mail', 'username', 'user name', 'account'];
                    const inputs = document.querySelectorAll('input[type="text"], input[type="email"], input:not([type])');
                    
                    for (const input of inputs) {
                        const rect = input.getBoundingClientRect();
                        if (rect.width === 0 || rect.height === 0) continue;
                        
                        // Check nearby text
                        const nearbyText = [];
                        const labels = document.querySelectorAll('label');
                        
                        for (const label of labels) {
                            if (label.getAttribute('for') === input.id || label.contains(input)) {
                                nearbyText.push(label.textContent.toLowerCase());
                            }
                        }
                        
                        // Check parent elements for text
                        let parent = input.parentElement;
                        let depth = 0;
                        while (parent && depth < 3) {
                            nearbyText.push(parent.textContent.toLowerCase());
                            parent = parent.parentElement;
                            depth++;
                        }
                        
                        // Check if any nearby text contains email keywords
                        const hasEmailKeyword = emailKeywords.some(keyword => 
                            nearbyText.some(text => text.includes(keyword))
                        );
                        
                        if (hasEmailKeyword) {
                            candidates.push({
                                selector: `input[name="${input.name}"]` || 
                                         `input[id="${input.id}"]` || 
                                         `input[type="${input.type}"]`,
                                score: 1
                            });
                        }
                    }
                    
                    return candidates;
                }
            """)
            
            if email_candidates:
                best_candidate = email_candidates[0]
                try:
                    element = page.locator(best_candidate['selector']).first
                    if await element.is_visible():
                        elements['email'] = element
                except:
                    pass
            
        except Exception as e:
            self.logger.debug(f"Fuzzy email detection failed: {str(e)}")
        
        # Similar fuzzy detection for password fields
        try:
            password_candidates = await page.evaluate("""
                () => {
                    const candidates = [];
                    const passwordKeywords = ['password', 'pass', 'pwd', 'secret'];
                    const inputs = document.querySelectorAll('input[type="password"]');
                    
                    if (inputs.length > 0) {
                        return [{selector: 'input[type="password"]', score: 1}];
                    }
                    
                    return candidates;
                }
            """)
            
            if password_candidates:
                try:
                    element = page.locator('input[type="password"]').first
                    if await element.is_visible():
                        elements['password'] = element
                except:
                    pass
                    
        except Exception as e:
            self.logger.debug(f"Fuzzy password detection failed: {str(e)}")
        
        return elements
    
    async def _last_resort_detection(self, page: Page) -> Dict[str, Any]:
        """Last resort: very broad selectors"""
        elements = {}
        
        try:
            # Just find any visible inputs
            all_inputs = page.locator('input:visible')
            input_count = await all_inputs.count()
            
            if input_count >= 2:
                # Assume first text input is email, first password is password
                for i in range(input_count):
                    input_elem = all_inputs.nth(i)
                    input_type = await input_elem.get_attribute('type') or 'text'
                    
                    if input_type.lower() == 'password' and 'password' not in elements:
                        elements['password'] = input_elem
                    elif (input_type.lower() in ['text', 'email'] and 
                          'email' not in elements):
                        elements['email'] = input_elem
            
            # Find any button
            buttons = page.locator('button:visible, input[type="submit"]:visible')
            if await buttons.count() > 0:
                elements['submit'] = buttons.first
            
        except Exception as e:
            self.logger.debug(f"Last resort detection failed: {str(e)}")
        
        return elements
    
    def _validate_elements(self, elements: Dict[str, Any]) -> bool:
        """
        Validate that we have minimum required elements
        
        Args:
            elements: Dictionary of detected elements
            
        Returns:
            bool: True if we have sufficient elements
        """
        # We need at least email/username and password
        required_elements = ['email', 'password']
        
        for required in required_elements:
            if required not in elements:
                return False
        
        # Submit button is helpful but not strictly required
        # (we can press Enter on password field)
        
        return True