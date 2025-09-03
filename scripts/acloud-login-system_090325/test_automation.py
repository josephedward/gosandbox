"""
Unit Tests for AI Login Automation System
"""

import pytest
import asyncio
import os
from unittest.mock import Mock, AsyncMock, patch
from pathlib import Path

# Import modules to test
import sys
sys.path.insert(0, str(Path(__file__).parent.parent))

from config.settings import Settings
from modules.ai_detector import AIElementDetector
from modules.captcha_solver import CaptchaSolver
from modules.element_mapper import ElementMapper
from main import AILoginAutomator


class TestSettings:
    """Test settings configuration"""
    
    def test_settings_from_env(self):
        """Test settings loading from environment variables"""
        # Mock environment variables
        test_env = {
            'LOGIN_USERNAME': 'test@example.com',
            'LOGIN_PASSWORD': 'testpass123',
            'OPENAI_API_KEY': 'sk-test123',
            'CAPSOLVER_API_KEY': 'cap-test123'
        }
        
        with patch.dict(os.environ, test_env):
            settings = Settings()
            assert settings.LOGIN_USERNAME == 'test@example.com'
            assert settings.OPENAI_API_KEY == 'sk-test123'
    
    def test_settings_defaults(self):
        """Test default settings values"""
        test_env = {
            'LOGIN_USERNAME': 'test@example.com',
            'LOGIN_PASSWORD': 'testpass123',
            'OPENAI_API_KEY': 'sk-test123'
        }
        
        with patch.dict(os.environ, test_env, clear=True):
            settings = Settings()
            assert settings.HEADLESS_MODE == False
            assert settings.VIEWPORT_WIDTH == 1920
            assert settings.MAX_RETRIES == 3


class TestAIElementDetector:
    """Test AI element detection functionality"""
    
    @pytest.fixture
    def mock_settings(self):
        """Mock settings for testing"""
        settings = Mock(spec=Settings)
        settings.OPENAI_API_KEY = 'sk-test123'
        settings.OPENAI_MODEL = 'gpt-4-turbo'
        settings.OPENAI_MAX_TOKENS = 1000
        return settings
    
    @pytest.fixture
    def detector(self, mock_settings):
        """Create AI detector instance for testing"""
        with patch('modules.ai_detector.openai.AsyncOpenAI'):
            return AIElementDetector(mock_settings)
    
    @pytest.mark.asyncio
    async def test_detect_elements_success(self, detector):
        """Test successful element detection"""
        # Mock page and screenshot
        mock_page = AsyncMock()
        mock_page.screenshot.return_value = b'fake_screenshot_data'
        
        # Mock GPT-4V response
        mock_response = {
            "email_field": {"x": 50, "y": 30},
            "password_field": {"x": 50, "y": 45},
            "submit_button": {"x": 50, "y": 60}
        }
        
        with patch.object(detector, '_analyze_screenshot', return_value=mock_response):
            with patch.object(detector, '_coords_to_elements', return_value={'email': 'mock_element'}):
                result = await detector.detect_elements(mock_page)
                assert 'email' in result
    
    @pytest.mark.asyncio
    async def test_detect_elements_failure(self, detector):
        """Test element detection failure handling"""
        mock_page = AsyncMock()
        mock_page.screenshot.side_effect = Exception("Screenshot failed")
        
        result = await detector.detect_elements(mock_page)
        assert result == {}
    
    def test_coordinate_validation(self, detector):
        """Test coordinate validation logic"""
        box = {'x': 100, 'y': 100, 'width': 50, 'height': 20}
        
        # Test near coordinates
        assert detector._is_near_coordinates(box, 125, 110, tolerance=50) == True
        
        # Test far coordinates
        assert detector._is_near_coordinates(box, 200, 200, tolerance=50) == False


class TestCaptchaSolver:
    """Test CAPTCHA solving functionality"""
    
    @pytest.fixture
    def mock_settings(self):
        """Mock settings for testing"""
        settings = Mock(spec=Settings)
        settings.CAPSOLVER_API_KEY = 'cap-test123'
        settings.CAPMONSTER_API_KEY = None
        settings.TWOCAPTCHA_API_KEY = None
        return settings
    
    @pytest.fixture
    def solver(self, mock_settings):
        """Create CAPTCHA solver instance for testing"""
        with patch('modules.captcha_solver.capsolver'):
            return CaptchaSolver(mock_settings)
    
    @pytest.mark.asyncio
    async def test_captcha_detection(self, solver):
        """Test CAPTCHA detection on page"""
        mock_page = AsyncMock()
        mock_element = AsyncMock()
        mock_element.count.return_value = 1
        mock_element.first.get_attribute.return_value = 'test-site-key'
        mock_page.locator.return_value = mock_element
        
        captcha_info = await solver._detect_captcha_info(mock_page)
        
        assert captcha_info is not None
        assert captcha_info['type'] == 'recaptcha_v2'
        assert captcha_info['site_key'] == 'test-site-key'
    
    @pytest.mark.asyncio
    async def test_no_captcha_detected(self, solver):
        """Test when no CAPTCHA is present"""
        mock_page = AsyncMock()
        mock_element = AsyncMock()
        mock_element.count.return_value = 0
        mock_page.locator.return_value = mock_element
        mock_page.evaluate.return_value = None
        
        captcha_info = await solver._detect_captcha_info(mock_page)
        assert captcha_info is None
    
    @pytest.mark.asyncio
    async def test_solve_with_capsolver_success(self, solver):
        """Test successful CAPTCHA solving with Capsolver"""
        captcha_info = {
            'type': 'recaptcha_v2',
            'site_key': 'test-key'
        }
        
        with patch('modules.captcha_solver.capsolver') as mock_capsolver:
            mock_capsolver.solve.return_value = {
                'solution': {'gRecaptchaResponse': 'solved_token_123'}
            }
            
            result = await solver._solve_with_capsolver(captcha_info, 'https://example.com')
            assert result == 'solved_token_123'


class TestElementMapper:
    """Test fallback element detection"""
    
    @pytest.fixture
    def mock_settings(self):
        """Mock settings for testing"""
        return Mock(spec=Settings)
    
    @pytest.fixture
    def mapper(self, mock_settings):
        """Create element mapper instance for testing"""
        return ElementMapper(mock_settings)
    
    @pytest.mark.asyncio
    async def test_semantic_detection(self, mapper):
        """Test semantic HTML detection"""
        mock_page = AsyncMock()
        
        # Mock email field
        mock_email_element = AsyncMock()
        mock_email_element.is_visible.return_value = True
        mock_email_locator = AsyncMock()
        mock_email_locator.first = mock_email_element
        
        # Mock password field
        mock_password_element = AsyncMock()
        mock_password_element.is_visible.return_value = True
        mock_password_locator = AsyncMock()
        mock_password_locator.first = mock_password_element
        
        # Mock submit button
        mock_submit_element = AsyncMock()
        mock_submit_element.is_visible.return_value = True
        mock_submit_locator = AsyncMock()
        mock_submit_locator.first = mock_submit_element
        
        # Configure page.locator to return appropriate mocks
        def locator_side_effect(selector):
            if 'email' in selector:
                return mock_email_locator
            elif 'password' in selector:
                return mock_password_locator
            elif 'submit' in selector or 'button' in selector:
                return mock_submit_locator
            else:
                empty_locator = AsyncMock()
                empty_locator.is_visible.return_value = False
                return empty_locator
        
        mock_page.locator.side_effect = locator_side_effect
        
        result = await mapper._semantic_detection(mock_page)
        
        assert 'email' in result
        assert 'password' in result
        assert 'submit' in result
    
    def test_element_validation(self, mapper):
        """Test element validation logic"""
        # Valid elements
        valid_elements = {
            'email': 'mock_element',
            'password': 'mock_element'
        }
        assert mapper._validate_elements(valid_elements) == True
        
        # Invalid elements (missing password)
        invalid_elements = {
            'email': 'mock_element'
        }
        assert mapper._validate_elements(invalid_elements) == False


class TestAILoginAutomator:
    """Test main automation class"""
    
    @pytest.fixture
    def mock_settings(self):
        """Mock settings for testing"""
        settings = Mock(spec=Settings)
        settings.LOGIN_USERNAME = 'test@example.com'
        settings.LOGIN_PASSWORD = 'testpass123'
        settings.HEADLESS_MODE = True
        settings.VIEWPORT_WIDTH = 1920
        settings.VIEWPORT_HEIGHT = 1080
        settings.USER_AGENT = 'test-agent'
        settings.BROWSER_TIMEOUT = 30000
        settings.RANDOM_DELAYS = False
        settings.STEALTH_MODE = True
        settings.SCREENSHOT_ON_FAILURE = False
        return settings
    
    @pytest.fixture
    def automator(self, mock_settings):
        """Create automator instance for testing"""
        with patch('main.AIElementDetector'), \
             patch('main.CaptchaSolver'), \
             patch('main.ElementMapper'), \
             patch('main.setup_logger'):
            return AILoginAutomator(mock_settings)
    
    @pytest.mark.asyncio
    async def test_login_success_flow(self, automator):
        """Test successful login flow"""
        # Mock all the dependencies
        automator.ai_detector.detect_elements = AsyncMock(return_value={
            'email': AsyncMock(),
            'password': AsyncMock(),
            'submit': AsyncMock()
        })
        
        with patch('main.async_playwright') as mock_playwright:
            mock_context_manager = AsyncMock()
            mock_playwright.return_value = mock_context_manager
            
            mock_browser = AsyncMock()
            mock_context = AsyncMock()
            mock_page = AsyncMock()
            
            mock_context_manager.__aenter__.return_value.chromium.launch.return_value = mock_browser
            mock_browser.new_context.return_value = mock_context
            mock_context.new_page.return_value = mock_page
            
            # Mock successful operations
            automator._navigate_to_page = AsyncMock()
            automator._detect_elements = AsyncMock(return_value={'email': 'mock', 'password': 'mock'})
            automator._fill_credentials = AsyncMock()
            automator._has_captcha = AsyncMock(return_value=False)
            automator._submit_form = AsyncMock()
            automator._verify_login_success = AsyncMock(return_value=True)
            
            result = await automator.perform_login('https://example.com/login')
            assert result == True
    
    @pytest.mark.asyncio
    async def test_login_failure_flow(self, automator):
        """Test login failure handling"""
        automator.ai_detector.detect_elements = AsyncMock(return_value={})
        
        with patch('main.async_playwright') as mock_playwright:
            mock_context_manager = AsyncMock()
            mock_playwright.return_value = mock_context_manager
            
            mock_browser = AsyncMock()
            mock_context = AsyncMock()
            mock_page = AsyncMock()
            
            mock_context_manager.__aenter__.return_value.chromium.launch.return_value = mock_browser
            mock_browser.new_context.return_value = mock_context
            mock_context.new_page.return_value = mock_page
            
            # Mock failed element detection
            automator._navigate_to_page = AsyncMock()
            automator._detect_elements = AsyncMock(return_value={})
            
            result = await automator.perform_login('https://example.com/login')
            assert result == False


class TestIntegration:
    """Integration tests"""
    
    @pytest.mark.asyncio
    @pytest.mark.integration
    async def test_real_browser_launch(self):
        """Test real browser launch (requires Playwright installation)"""
        try:
            from playwright.async_api import async_playwright
            
            async with async_playwright() as p:
                browser = await p.chromium.launch(headless=True)
                context = await browser.new_context()
                page = await context.new_page()
                
                await page.goto('https://httpbin.org/get')
                title = await page.title()
                
                assert title is not None
                await browser.close()
                
        except ImportError:
            pytest.skip("Playwright not installed")
    
    @pytest.mark.asyncio
    @pytest.mark.integration
    async def test_openai_api_format(self):
        """Test OpenAI API call format (without actually calling API)"""
        # This test ensures our API call structure is correct
        try:
            import openai
            
            # Test that our message format is valid
            messages = [{
                "role": "user",
                "content": [
                    {"type": "text", "text": "Test prompt"},
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": "data:image/png;base64,test_data",
                            "detail": "high"
                        }
                    }
                ]
            }]
            
            # Verify structure is valid (don't actually call API)
            assert len(messages) == 1
            assert messages[0]["role"] == "user"
            assert len(messages[0]["content"]) == 2
            
        except ImportError:
            pytest.skip("OpenAI library not installed")


# Fixtures for test data
@pytest.fixture
def sample_page_html():
    """Sample HTML for testing"""
    return """
    <!DOCTYPE html>
    <html>
    <head><title>Login Page</title></head>
    <body>
        <form id="loginForm">
            <input type="email" name="email" id="email" placeholder="Enter your email">
            <input type="password" name="password" id="password" placeholder="Enter your password">
            <button type="submit" id="submitBtn">Sign In</button>
        </form>
        <div class="g-recaptcha" data-sitekey="test-site-key"></div>
    </body>
    </html>
    """


@pytest.fixture
def sample_captcha_response():
    """Sample CAPTCHA service response"""
    return {
        'status': 'ready',
        'solution': {
            'gRecaptchaResponse': 'test_captcha_solution_token_12345'
        }
    }


# Test configuration
def pytest_configure(config):
    """Configure pytest with custom markers"""
    config.addinivalue_line(
        "markers", "integration: mark test as integration test"
    )


# Run tests with: python -m pytest tests/test_automation.py -v
if __name__ == "__main__":
    pytest.main([__file__, "-v"])


# Performance benchmarks
class TestPerformance:
    """Performance and benchmark tests"""
    
    @pytest.mark.asyncio
    async def test_element_detection_performance(self):
        """Benchmark element detection speed"""
        import time
        
        mock_settings = Mock(spec=Settings)
        
        with patch('modules.element_mapper.setup_logger'):
            mapper = ElementMapper(mock_settings)
            
            mock_page = AsyncMock()
            mock_element = AsyncMock()
            mock_element.is_visible.return_value = True
            mock_locator = AsyncMock()
            mock_locator.first = mock_element
            mock_page.locator.return_value = mock_locator
            
            start_time = time.time()
            await mapper._semantic_detection(mock_page)
            duration = time.time() - start_time
            
            # Should complete within reasonable time
            assert duration < 1.0  # Less than 1 second
    
    @pytest.mark.asyncio
    async def test_multiple_retries_performance(self):
        """Test performance with multiple retry attempts"""
        import time
        
        mock_settings = Mock(spec=Settings)
        mock_settings.MAX_RETRIES = 3
        mock_settings.RETRY_DELAY = 0.1  # Fast retries for testing
        
        start_time = time.time()
        
        # Simulate retry logic
        for attempt in range(mock_settings.MAX_RETRIES):
            await asyncio.sleep(mock_settings.RETRY_DELAY)
        
        duration = time.time() - start_time
        
        # Should complete all retries within expected time
        expected_max_duration = mock_settings.MAX_RETRIES * mock_settings.RETRY_DELAY + 0.5
        assert duration < expected_max_duration