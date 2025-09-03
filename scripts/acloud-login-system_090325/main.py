#!/usr/bin/env python3
"""
AI-Powered Login Automation System
Main entry point for the automation system
"""

import asyncio
import os
import sys
from pathlib import Path
from typing import Optional, Dict, Any
import traceback

# Add project root to Python path
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

from modules.ai_detector import AIElementDetector
from modules.captcha_solver import CaptchaSolver
from modules.element_mapper import ElementMapper
from modules.logger import setup_logger
from config.settings import Settings


class AILoginAutomator:
    """Main automation class that orchestrates the login process"""
    
    def __init__(self, config: Optional[Settings] = None):
        self.config = config or Settings()
        self.logger = setup_logger()
        self.ai_detector = AIElementDetector(self.config)
        self.captcha_solver = CaptchaSolver(self.config)
        self.element_mapper = ElementMapper(self.config)
        
    async def perform_login(self, url: str) -> bool:
        """
        Main login automation function
        
        Args:
            url: Target login URL
            
        Returns:
            bool: True if login successful, False otherwise
        """
        self.logger.info(f"Starting login automation for: {url}")
        
        try:
            from playwright.async_api import async_playwright
            
            async with async_playwright() as p:
                # Launch browser with stealth settings
                browser = await self._launch_browser(p)
                context = await self._create_context(browser)
                page = await context.new_page()
                
                try:
                    # Navigate to login page
                    await self._navigate_to_page(page, url)
                    
                    # Detect form elements using AI
                    elements = await self._detect_elements(page)
                    
                    if not elements:
                        self.logger.error("Failed to detect form elements")
                        return False
                    
                    # Fill login credentials
                    await self._fill_credentials(page, elements)
                    
                    # Handle CAPTCHA if present
                    if await self._has_captcha(page):
                        captcha_solved = await self._solve_captcha(page)
                        if not captcha_solved:
                            self.logger.warning("CAPTCHA solving failed, proceeding anyway")
                    
                    # Submit form
                    await self._submit_form(page, elements)
                    
                    # Verify login success
                    success = await self._verify_login_success(page)
                    
                    if success:
                        self.logger.info("Login automation completed successfully")
                    else:
                        self.logger.error("Login automation failed")
                    
                    return success
                    
                except Exception as e:
                    self.logger.error(f"Login process failed: {str(e)}")
                    if self.config.SCREENSHOT_ON_FAILURE:
                        await self._take_failure_screenshot(page)
                    return False
                    
                finally:
                    await browser.close()
                    
        except Exception as e:
            self.logger.error(f"Browser setup failed: {str(e)}")
            self.logger.debug(traceback.format_exc())
            return False
    
    async def _launch_browser(self, playwright):
        """Launch browser with optimal settings"""
        browser_options = {
            'headless': self.config.HEADLESS_MODE,
            'args': [
                '--no-sandbox',
                '--disable-blink-features=AutomationControlled',
                '--disable-dev-shm-usage',
                '--disable-gpu',
                '--no-first-run',
                '--no-default-browser-check',
                '--disable-default-apps',
                '--disable-popup-blocking'
            ]
        }
        
        if self.config.ENABLE_PROXY and self.config.PROXY_URL:
            proxy_parts = self.config.PROXY_URL.replace('http://', '').split('@')
            if len(proxy_parts) == 2:
                auth_parts = proxy_parts[0].split(':')
                server_parts = proxy_parts[1].split(':')
                browser_options['proxy'] = {
                    'server': f"http://{proxy_parts[1]}",
                    'username': auth_parts[0],
                    'password': auth_parts[1]
                }
        
        return await playwright.chromium.launch(**browser_options)
    
    async def _create_context(self, browser):
        """Create browser context with stealth settings"""
        context_options = {
            'viewport': {
                'width': self.config.VIEWPORT_WIDTH,
                'height': self.config.VIEWPORT_HEIGHT
            },
            'user_agent': self.config.USER_AGENT,
            'locale': 'en-US',
            'timezone_id': 'America/New_York',
            'permissions': ['geolocation', 'notifications'],
            'extra_http_headers': {
                'Accept-Language': 'en-US,en;q=0.9',
                'Accept-Encoding': 'gzip, deflate, br',
                'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
                'Connection': 'keep-alive',
                'Upgrade-Insecure-Requests': '1',
            }
        }
        
        context = await browser.new_context(**context_options)
        
        # Add stealth scripts
        if self.config.STEALTH_MODE:
            await context.add_init_script("""
                Object.defineProperty(navigator, 'webdriver', {
                    get: () => undefined,
                });
                
                Object.defineProperty(navigator, 'plugins', {
                    get: () => [1, 2, 3, 4, 5],
                });
                
                Object.defineProperty(navigator, 'languages', {
                    get: () => ['en-US', 'en'],
                });
                
                window.chrome = {
                    runtime: {},
                };
            """)
        
        return context
    
    async def _navigate_to_page(self, page, url: str):
        """Navigate to the target page with proper waiting"""
        self.logger.info(f"Navigating to: {url}")
        
        await page.goto(url, wait_until='networkidle', timeout=self.config.BROWSER_TIMEOUT)
        
        # Random delay for human-like behavior
        if self.config.RANDOM_DELAYS:
            import random
            delay = random.uniform(1, 3)
            self.logger.debug(f"Adding random delay: {delay:.2f}s")
            await asyncio.sleep(delay)
    
    async def _detect_elements(self, page) -> Dict[str, Any]:
        """Detect form elements using AI and fallback methods"""
        self.logger.info("Starting element detection")
        
        # Try AI detection first
        try:
            elements = await self.ai_detector.detect_elements(page)
            if elements:
                self.logger.info("AI detection successful")
                return elements
        except Exception as e:
            self.logger.warning(f"AI detection failed: {str(e)}")
        
        # Fallback to traditional methods
        self.logger.info("Falling back to traditional element detection")
        elements = await self.element_mapper.detect_elements_fallback(page)
        
        return elements
    
    async def _fill_credentials(self, page, elements: Dict[str, Any]):
        """Fill login credentials in detected form elements"""
        self.logger.info("Filling login credentials")
        
        # Fill username/email
        if 'email' in elements:
            email_element = elements['email']
            await email_element.clear()
            await email_element.fill(self.config.LOGIN_USERNAME)
            self.logger.debug("Email field filled")
            
            if self.config.RANDOM_DELAYS:
                await asyncio.sleep(0.5)
        
        # Fill password
        if 'password' in elements:
            password_element = elements['password']
            await password_element.clear()
            await password_element.fill(self.config.LOGIN_PASSWORD)
            self.logger.debug("Password field filled")
            
            if self.config.RANDOM_DELAYS:
                await asyncio.sleep(0.5)
    
    async def _has_captcha(self, page) -> bool:
        """Check if CAPTCHA is present on the page"""
        captcha_selectors = [
            '.g-recaptcha',
            '[data-sitekey]',
            '#recaptcha',
            '.recaptcha',
            '.captcha',
            'iframe[src*="recaptcha"]',
            '.h-captcha',
            '[data-hcaptcha-sitekey]'
        ]
        
        for selector in captcha_selectors:
            try:
                element = page.locator(selector)
                if await element.count() > 0:
                    self.logger.info(f"CAPTCHA detected: {selector}")
                    return True
            except:
                continue
        
        return False
    
    async def _solve_captcha(self, page) -> bool:
        """Solve CAPTCHA using external service"""
        self.logger.info("Attempting to solve CAPTCHA")
        
        try:
            success = await self.captcha_solver.solve_captcha(page)
            if success:
                self.logger.info("CAPTCHA solved successfully")
            else:
                self.logger.error("CAPTCHA solving failed")
            return success
        except Exception as e:
            self.logger.error(f"CAPTCHA solving error: {str(e)}")
            return False
    
    async def _submit_form(self, page, elements: Dict[str, Any]):
        """Submit the login form"""
        self.logger.info("Submitting login form")
        
        if 'submit' in elements:
            submit_element = elements['submit']
            await submit_element.click()
            self.logger.debug("Submit button clicked")
        else:
            # Try pressing Enter on password field
            if 'password' in elements:
                await elements['password'].press('Enter')
                self.logger.debug("Enter pressed on password field")
        
        # Wait for navigation or response
        try:
            await page.wait_for_load_state('networkidle', timeout=10000)
        except:
            # Continue if timeout occurs
            pass
    
    async def _verify_login_success(self, page) -> bool:
        """Verify if login was successful"""
        self.logger.info("Verifying login success")
        
        # Wait a bit for page to load
        await asyncio.sleep(2)
        
        # Success indicators
        success_indicators = [
            # Common success patterns
            ':has-text("dashboard")',
            ':has-text("profile")',
            ':has-text("logout")',
            ':has-text("welcome")',
            ':has-text("account")',
            '.user-menu',
            '.user-avatar',
            '[data-testid*="user"]',
            '[aria-label*="profile"]',
            '.navbar-user',
            '.user-dropdown'
        ]
        
        # Check for success indicators
        for indicator in success_indicators:
            try:
                element = page.locator(indicator)
                if await element.count() > 0:
                    self.logger.info(f"Success indicator found: {indicator}")
                    return True
            except:
                continue
        
        # Check if still on login page (failure indicator)
        failure_indicators = [
            'input[type="password"]',
            ':has-text("sign in")',
            ':has-text("login")',
            ':has-text("email")',
            ':has-text("password")',
            '.login-form',
            '.signin-form'
        ]
        
        for indicator in failure_indicators:
            try:
                element = page.locator(indicator)
                if await element.count() > 0:
                    self.logger.warning(f"Still on login page: {indicator}")
                    return False
            except:
                continue
        
        # Check URL change
        current_url = page.url
        if 'login' not in current_url.lower() and 'signin' not in current_url.lower():
            self.logger.info("URL changed from login page")
            return True
        
        self.logger.warning("Login success verification inconclusive")
        return False
    
    async def _take_failure_screenshot(self, page):
        """Take screenshot on failure for debugging"""
        try:
            screenshot_path = f"failure_screenshot_{int(asyncio.get_event_loop().time())}.png"
            await page.screenshot(path=screenshot_path, full_page=True)
            self.logger.info(f"Failure screenshot saved: {screenshot_path}")
        except Exception as e:
            self.logger.error(f"Failed to take screenshot: {str(e)}")


async def main():
    """Main entry point"""
    
    # Load configuration
    config = Settings()
    
    # Initialize automator
    automator = AILoginAutomator(config)
    
    # Target URL - can be customized
    target_url = "https://learn.acloud.guru/cloud-playground/cloud-sandboxes"
    
    # Perform login
    success = await automator.perform_login(target_url)
    
    if success:
        print("✅ Login automation completed successfully!")
        return 0
    else:
        print("❌ Login automation failed!")
        return 1


if __name__ == "__main__":
    try:
        exit_code = asyncio.run(main())
        sys.exit(exit_code)
    except KeyboardInterrupt:
        print("\n🛑 Login automation interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"💥 Unexpected error: {str(e)}")
        sys.exit(1)