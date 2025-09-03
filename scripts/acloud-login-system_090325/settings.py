"""
Settings Configuration Module
Manages all configuration settings from environment variables
"""

import os
from typing import Optional
from dotenv import load_dotenv
from pydantic import BaseSettings, Field

# Load environment variables
load_dotenv()


class Settings(BaseSettings):
    """Application settings loaded from environment variables"""
    
    # Authentication Credentials
    LOGIN_USERNAME: str = Field(..., env="LOGIN_USERNAME", description="Login username/email")
    LOGIN_PASSWORD: str = Field(..., env="LOGIN_PASSWORD", description="Login password")
    
    # OpenAI Configuration
    OPENAI_API_KEY: str = Field(..., env="OPENAI_API_KEY", description="OpenAI API key")
    OPENAI_MODEL: str = Field(default="gpt-4-turbo", env="OPENAI_MODEL", description="OpenAI model to use")
    OPENAI_MAX_TOKENS: int = Field(default=1000, env="OPENAI_MAX_TOKENS", description="Max tokens for OpenAI responses")
    
    # CAPTCHA Service API Keys
    CAPSOLVER_API_KEY: Optional[str] = Field(default=None, env="CAPSOLVER_API_KEY", description="Capsolver API key")
    CAPMONSTER_API_KEY: Optional[str] = Field(default=None, env="CAPMONSTER_API_KEY", description="CapMonster API key")
    TWOCAPTCHA_API_KEY: Optional[str] = Field(default=None, env="TWOCAPTCHA_API_KEY", description="2Captcha API key")
    
    # Browser Configuration
    HEADLESS_MODE: bool = Field(default=False, env="HEADLESS_MODE", description="Run browser in headless mode")
    BROWSER_TIMEOUT: int = Field(default=30000, env="BROWSER_TIMEOUT", description="Browser timeout in milliseconds")
    VIEWPORT_WIDTH: int = Field(default=1920, env="VIEWPORT_WIDTH", description="Browser viewport width")
    VIEWPORT_HEIGHT: int = Field(default=1080, env="VIEWPORT_HEIGHT", description="Browser viewport height")
    USER_AGENT: str = Field(
        default="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
        env="USER_AGENT",
        description="Browser user agent string"
    )
    
    # Automation Settings
    MAX_RETRIES: int = Field(default=3, env="MAX_RETRIES", description="Maximum retry attempts")
    RETRY_DELAY: float = Field(default=2.0, env="RETRY_DELAY", description="Delay between retries in seconds")
    SCREENSHOT_ON_FAILURE: bool = Field(default=True, env="SCREENSHOT_ON_FAILURE", description="Take screenshot on failure")
    ELEMENT_DETECTION_TIMEOUT: int = Field(default=10, env="ELEMENT_DETECTION_TIMEOUT", description="Element detection timeout in seconds")
    
    # Logging Configuration
    LOG_LEVEL: str = Field(default="INFO", env="LOG_LEVEL", description="Logging level")
    LOG_FILE: str = Field(default="automation.log", env="LOG_FILE", description="Log file name")
    ENABLE_FILE_LOGGING: bool = Field(default=True, env="ENABLE_FILE_LOGGING", description="Enable file logging")
    ENABLE_CONSOLE_LOGGING: bool = Field(default=True, env="ENABLE_CONSOLE_LOGGING", description="Enable console logging")
    
    # Security Settings
    ENABLE_PROXY: bool = Field(default=False, env="ENABLE_PROXY", description="Enable proxy usage")
    PROXY_URL: Optional[str] = Field(default=None, env="PROXY_URL", description="Proxy URL")
    STEALTH_MODE: bool = Field(default=True, env="STEALTH_MODE", description="Enable stealth mode")
    RANDOM_DELAYS: bool = Field(default=True, env="RANDOM_DELAYS", description="Enable random delays")
    
    class Config:
        env_file = ".env"
        case_sensitive = True
        
    def validate_required_keys(self):
        """Validate that required API keys are present"""
        if not self.OPENAI_API_KEY:
            raise ValueError("OPENAI_API_KEY is required")
        
        if not self.LOGIN_USERNAME or not self.LOGIN_PASSWORD:
            raise ValueError("LOGIN_USERNAME and LOGIN_PASSWORD are required")
        
        # Check that at least one CAPTCHA service is configured
        captcha_services = [
            self.CAPSOLVER_API_KEY,
            self.CAPMONSTER_API_KEY,
            self.TWOCAPTCHA_API_KEY
        ]
        
        if not any(captcha_services):
            print("⚠️  Warning: No CAPTCHA solving services configured. CAPTCHA solving will be disabled.")
    
    def get_captcha_services(self) -> list:
        """Get list of configured CAPTCHA services"""
        services = []
        
        if self.CAPSOLVER_API_KEY:
            services.append("Capsolver")
        if self.CAPMONSTER_API_KEY:
            services.append("CapMonster")
        if self.TWOCAPTCHA_API_KEY:
            services.append("2Captcha")
        
        return services