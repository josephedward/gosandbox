"""
CAPTCHA Solver Module
Integrates with multiple CAPTCHA solving services
"""

import asyncio
import re
from typing import Optional, Dict, Any
from playwright.async_api import Page

try:
    import capsolver
except ImportError:
    capsolver = None

from config.settings import Settings
from modules.logger import setup_logger


class CaptchaSolver:
    """Multi-service CAPTCHA solver with fallback support"""
    
    def __init__(self, config: Settings):
        self.config = config
        self.logger = setup_logger()
        self._init_services()
    
    def _init_services(self):
        """Initialize CAPTCHA solving services"""
        self.services = []
        
        # Capsolver
        if self.config.CAPSOLVER_API_KEY and capsolver:
            capsolver.api_key = self.config.CAPSOLVER_API_KEY
            self.services.append(('capsolver', self._solve_with_capsolver))
            self.logger.info("Capsolver service initialized")
        
        # CapMonster (via requests)
        if self.config.CAPMONSTER_API_KEY:
            self.services.append(('capmonster', self._solve_with_capmonster))
            self.logger.info("CapMonster service initialized")
        
        # 2Captcha (via requests)
        if self.config.TWOCAPTCHA_API_KEY:
            self.services.append(('2captcha', self._solve_with_2captcha))
            self.logger.info("2Captcha service initialized")
        
        if not self.services:
            self.logger.warning("No CAPTCHA solving services configured")
    
    async def solve_captcha(self, page: Page) -> bool:
        """
        Solve CAPTCHA on the current page
        
        Args:
            page: Playwright page object
            
        Returns:
            bool: True if CAPTCHA solved successfully
        """
        try:
            # Detect CAPTCHA type and get details
            captcha_info = await self._detect_captcha_info(page)
            
            if not captcha_info:
                self.logger.warning("No CAPTCHA detected or unable to extract info")
                return False
            
            self.logger.info(f"Detected {captcha_info['type']} CAPTCHA")
            
            # Try each service until one succeeds
            for service_name, solve_func in self.services:
                try:
                    self.logger.info(f"Attempting to solve with {service_name}")
                    solution = await solve_func(captcha_info, page.url)
                    
                    if solution:
                        # Inject solution into page
                        success = await self._inject_solution(page, captcha_info, solution)
                        if success:
                            self.logger.info(f"CAPTCHA solved successfully with {service_name}")
                            return True
                        else:
                            self.logger.warning(f"Failed to inject solution from {service_name}")
                    
                except Exception as e:
                    self.logger.warning(f"Service {service_name} failed: {str(e)}")
                    continue
            
            self.logger.error("All CAPTCHA solving services failed")
            return False
            
        except Exception as e:
            self.logger.error(f"CAPTCHA solving error: {str(e)}")
            return False
    
    async def _detect_captcha_info(self, page: Page) -> Optional[Dict[str, Any]]:
        """
        Detect CAPTCHA type and extract necessary information
        
        Args:
            page: Playwright page object
            
        Returns:
            Dict with CAPTCHA information or None
        """
        # Check for reCAPTCHA v2
        recaptcha_v2 = page.locator('.g-recaptcha, [data-sitekey]')
        if await recaptcha_v2.count() > 0:
            try:
                site_key = await recaptcha_v2.first.get_attribute('data-sitekey')
                if site_key:
                    return {
                        'type': 'recaptcha_v2',
                        'site_key': site_key,
                        'element': recaptcha_v2.first
                    }
            except:
                pass
        
        # Check for reCAPTCHA v3 (invisible)
        recaptcha_v3_scripts = await page.evaluate("""
            () => {
                const scripts = Array.from(document.querySelectorAll('script'));
                for (const script of scripts) {
                    const text = script.textContent || script.innerText || '';
                    const match = text.match(/grecaptcha\\.render\\(['"]([^'"]+)['"].*?['"]([A-Za-z0-9_-]+)['"].*?\\)/);
                    if (match) {
                        return { elementId: match[1], siteKey: match[2] };
                    }
                }
                return null;
            }
        """)
        
        if recaptcha_v3_scripts:
            return {
                'type': 'recaptcha_v3',
                'site_key': recaptcha_v3_scripts['siteKey'],
                'element_id': recaptcha_v3_scripts['elementId']
            }
        
        # Check for hCaptcha
        hcaptcha = page.locator('.h-captcha, [data-hcaptcha-sitekey]')
        if await hcaptcha.count() > 0:
            try:
                site_key = await hcaptcha.first.get_attribute('data-hcaptcha-sitekey')
                if site_key:
                    return {
                        'type': 'hcaptcha',
                        'site_key': site_key,
                        'element': hcaptcha.first
                    }
            except:
                pass
        
        # Check for other CAPTCHA types by looking for iframes
        captcha_iframes = await page.evaluate("""
            () => {
                const iframes = Array.from(document.querySelectorAll('iframe'));
                for (const iframe of iframes) {
                    const src = iframe.src || '';
                    if (src.includes('recaptcha')) {
                        return { type: 'recaptcha_iframe', src: src };
                    }
                    if (src.includes('hcaptcha')) {
                        return { type: 'hcaptcha_iframe', src: src };
                    }
                }
                return null;
            }
        """)
        
        if captcha_iframes:
            # Extract site key from iframe src
            site_key_match = re.search(r'[?&]k=([A-Za-z0-9_-]+)', captcha_iframes['src'])
            if site_key_match:
                captcha_type = 'recaptcha_v2' if 'recaptcha' in captcha_iframes['type'] else 'hcaptcha'
                return {
                    'type': captcha_type,
                    'site_key': site_key_match.group(1),
                    'iframe_src': captcha_iframes['src']
                }
        
        return None
    
    async def _solve_with_capsolver(self, captcha_info: Dict[str, Any], page_url: str) -> Optional[str]:
        """Solve CAPTCHA using Capsolver service"""
        if not capsolver:
            raise Exception("Capsolver library not installed")
        
        try:
            task_data = {
                "websiteURL": page_url,
                "websiteKey": captcha_info['site_key']
            }
            
            if captcha_info['type'] == 'recaptcha_v2':
                task_data["type"] = "ReCaptchaV2TaskProxyless"
            elif captcha_info['type'] == 'recaptcha_v3':
                task_data["type"] = "ReCaptchaV3TaskProxyless"
                task_data["pageAction"] = "submit"  # Default action
            elif captcha_info['type'] == 'hcaptcha':
                task_data["type"] = "HCaptchaTaskProxyless"
            else:
                raise Exception(f"Unsupported CAPTCHA type: {captcha_info['type']}")
            
            self.logger.debug(f"Capsolver task data: {task_data}")
            
            # Solve CAPTCHA
            solution = capsolver.solve(task_data)
            
            if solution and 'solution' in solution:
                if captcha_info['type'] == 'recaptcha_v3':
                    return solution['solution']['gRecaptchaResponse']
                elif captcha_info['type'] in ['recaptcha_v2', 'hcaptcha']:
                    return solution['solution']['gRecaptchaResponse']
            
            return None
            
        except Exception as e:
            self.logger.error(f"Capsolver error: {str(e)}")
            return None
    
    async def _solve_with_capmonster(self, captcha_info: Dict[str, Any], page_url: str) -> Optional[str]:
        """Solve CAPTCHA using CapMonster service"""
        import httpx
        
        try:
            # CapMonster API implementation
            create_task_url = "https://api.capmonster.cloud/createTask"
            get_result_url = "https://api.capmonster.cloud/getTaskResult"
            
            # Prepare task
            if captcha_info['type'] == 'recaptcha_v2':
                task = {
                    "type": "NoCaptchaTaskProxyless",
                    "websiteURL": page_url,
                    "websiteKey": captcha_info['site_key']
                }
            elif captcha_info['type'] == 'recaptcha_v3':
                task = {
                    "type": "RecaptchaV3TaskProxyless",
                    "websiteURL": page_url,
                    "websiteKey": captcha_info['site_key'],
                    "pageAction": "submit"
                }
            elif captcha_info['type'] == 'hcaptcha':
                task = {
                    "type": "HCaptchaTaskProxyless",
                    "websiteURL": page_url,
                    "websiteKey": captcha_info['site_key']
                }
            else:
                raise Exception(f"Unsupported CAPTCHA type: {captcha_info['type']}")
            
            # Create task
            async with httpx.AsyncClient() as client:
                create_response = await client.post(create_task_url, json={
                    "clientKey": self.config.CAPMONSTER_API_KEY,
                    "task": task
                })
                
                create_data = create_response.json()
                if create_data.get('errorId') != 0:
                    raise Exception(f"CapMonster create task error: {create_data.get('errorDescription')}")
                
                task_id = create_data['taskId']
                
                # Poll for result
                for _ in range(30):  # 30 attempts, 5 seconds each = 2.5 minutes max
                    await asyncio.sleep(5)
                    
                    result_response = await client.post(get_result_url, json={
                        "clientKey": self.config.CAPMONSTER_API_KEY,
                        "taskId": task_id
                    })
                    
                    result_data = result_response.json()
                    
                    if result_data.get('status') == 'ready':
                        return result_data['solution']['gRecaptchaResponse']
                    elif result_data.get('status') == 'processing':
                        continue
                    else:
                        raise Exception(f"CapMonster error: {result_data.get('errorDescription')}")
                
                raise Exception("CapMonster timeout")
                
        except Exception as e:
            self.logger.error(f"CapMonster error: {str(e)}")
            return None
    
    async def _solve_with_2captcha(self, captcha_info: Dict[str, Any], page_url: str) -> Optional[str]:
        """Solve CAPTCHA using 2Captcha service"""
        import httpx
        
        try:
            submit_url = "http://2captcha.com/in.php"
            result_url = "http://2captcha.com/res.php"
            
            # Prepare parameters
            params = {
                "key": self.config.TWOCAPTCHA_API_KEY,
                "method": "userrecaptcha",
                "googlekey": captcha_info['site_key'],
                "pageurl": page_url,
                "json": 1
            }
            
            if captcha_info['type'] == 'recaptcha_v3':
                params["version"] = "v3"
                params["action"] = "submit"
                params["min_score"] = "0.3"
            
            # Submit CAPTCHA
            async with httpx.AsyncClient() as client:
                submit_response = await client.post(submit_url, data=params)
                submit_data = submit_response.json()
                
                if submit_data.get('status') != 1:
                    raise Exception(f"2Captcha submit error: {submit_data.get('error_text')}")
                
                captcha_id = submit_data['request']
                
                # Poll for result
                for _ in range(30):  # 30 attempts, 5 seconds each = 2.5 minutes max
                    await asyncio.sleep(5)
                    
                    result_response = await client.get(result_url, params={
                        "key": self.config.TWOCAPTCHA_API_KEY,
                        "action": "get",
                        "id": captcha_id,
                        "json": 1
                    })
                    
                    result_data = result_response.json()
                    
                    if result_data.get('status') == 1:
                        return result_data['request']
                    elif result_data.get('error_text') == 'CAPCHA_NOT_READY':
                        continue
                    else:
                        raise Exception(f"2Captcha error: {result_data.get('error_text')}")
                
                raise Exception("2Captcha timeout")
                
        except Exception as e:
            self.logger.error(f"2Captcha error: {str(e)}")
            return None
    
    async def _inject_solution(self, page: Page, captcha_info: Dict[str, Any], solution: str) -> bool:
        """
        Inject CAPTCHA solution into the page
        
        Args:
            page: Playwright page object
            captcha_info: CAPTCHA information
            solution: Solved CAPTCHA token
            
        Returns:
            bool: True if injection successful
        """
        try:
            if captcha_info['type'] in ['recaptcha_v2', 'recaptcha_v3']:
                # Inject reCAPTCHA solution
                await page.evaluate(f"""
                    () => {{
                        // Set the response in the textarea
                        const textarea = document.getElementById('g-recaptcha-response') || 
                                       document.querySelector('[name="g-recaptcha-response"]');
                        if (textarea) {{
                            textarea.value = '{solution}';
                            textarea.innerHTML = '{solution}';
                            textarea.style.display = 'block';
                        }}
                        
                        // Override grecaptcha.getResponse
                        if (typeof grecaptcha !== 'undefined') {{
                            grecaptcha.getResponse = function() {{ return '{solution}'; }};
                        }}
                        
                        // Trigger callback if available
                        if (window.recaptchaCallback) {{
                            window.recaptchaCallback('{solution}');
                        }}
                        
                        return true;
                    }}
                """)
                
            elif captcha_info['type'] == 'hcaptcha':
                # Inject hCaptcha solution
                await page.evaluate(f"""
                    () => {{
                        // Set the response in the textarea
                        const textarea = document.querySelector('[name="h-captcha-response"]') ||
                                       document.getElementById('h-captcha-response');
                        if (textarea) {{
                            textarea.value = '{solution}';
                            textarea.innerHTML = '{solution}';
                        }}
                        
                        // Override hcaptcha.getResponse
                        if (typeof hcaptcha !== 'undefined') {{
                            hcaptcha.getResponse = function() {{ return '{solution}'; }};
                        }}
                        
                        return true;
                    }}
                """)
            
            # Wait a bit for the solution to be processed
            await asyncio.sleep(1)
            
            self.logger.info("CAPTCHA solution injected successfully")
            return True
            
        except Exception as e:
            self.logger.error(f"Failed to inject CAPTCHA solution: {str(e)}")
            return False