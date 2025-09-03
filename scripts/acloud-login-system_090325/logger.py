"""
Logging Module
Centralized logging configuration for the automation system
"""

import sys
import logging
from pathlib import Path
from typing import Optional
from loguru import logger
import os


def setup_logger(
    log_level: str = "INFO",
    log_file: Optional[str] = None,
    enable_file_logging: bool = True,
    enable_console_logging: bool = True
) -> logging.Logger:
    """
    Setup centralized logging for the application
    
    Args:
        log_level: Logging level (DEBUG, INFO, WARNING, ERROR)
        log_file: Log file path
        enable_file_logging: Enable file logging
        enable_console_logging: Enable console logging
        
    Returns:
        Configured logger instance
    """
    
    # Remove default logger
    logger.remove()
    
    # Configure log format
    log_format = (
        "<green>{time:YYYY-MM-DD HH:mm:ss}</green> | "
        "<level>{level: <8}</level> | "
        "<cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> | "
        "<level>{message}</level>"
    )
    
    # Console logging
    if enable_console_logging:
        logger.add(
            sys.stdout,
            format=log_format,
            level=log_level,
            colorize=True,
            backtrace=True,
            diagnose=True
        )
    
    # File logging
    if enable_file_logging:
        if not log_file:
            log_file = os.getenv('LOG_FILE', 'automation.log')
        
        # Create logs directory if it doesn't exist
        log_path = Path(log_file)
        log_path.parent.mkdir(parents=True, exist_ok=True)
        
        logger.add(
            log_file,
            format=log_format,
            level=log_level,
            rotation="10 MB",
            retention="7 days",
            compression="zip",
            backtrace=True,
            diagnose=True
        )
    
    # Create standard logging adapter for compatibility
    class LoguruHandler(logging.Handler):
        def emit(self, record):
            try:
                level = logger.level(record.levelname).name
            except ValueError:
                level = record.levelno
            
            frame, depth = logging.currentframe(), 2
            while frame.f_code.co_filename == logging.__file__:
                frame = frame.f_back
                depth += 1
            
            logger.opt(depth=depth, exception=record.exc_info).log(level, record.getMessage())
    
    # Setup standard logging to use loguru
    logging.basicConfig(handlers=[LoguruHandler()], level=0, force=True)
    
    # Disable some noisy loggers
    logging.getLogger('urllib3').setLevel(logging.WARNING)
    logging.getLogger('httpx').setLevel(logging.WARNING)
    logging.getLogger('httpcore').setLevel(logging.WARNING)
    logging.getLogger('openai').setLevel(logging.WARNING)
    
    return logger


def get_logger(name: str = None):
    """Get a logger instance"""
    if name:
        return logger.bind(name=name)
    return logger


# Performance logging decorator
def log_performance(func_name: str = None):
    """Decorator to log function performance"""
    def decorator(func):
        import time
        import functools
        
        @functools.wraps(func)
        async def async_wrapper(*args, **kwargs):
            start_time = time.time()
            name = func_name or func.__name__
            
            logger.debug(f"Starting {name}")
            try:
                result = await func(*args, **kwargs)
                duration = time.time() - start_time
                logger.info(f"Completed {name} in {duration:.2f}s")
                return result
            except Exception as e:
                duration = time.time() - start_time
                logger.error(f"Failed {name} after {duration:.2f}s: {str(e)}")
                raise
                
        @functools.wraps(func)
        def sync_wrapper(*args, **kwargs):
            start_time = time.time()
            name = func_name or func.__name__
            
            logger.debug(f"Starting {name}")
            try:
                result = func(*args, **kwargs)
                duration = time.time() - start_time
                logger.info(f"Completed {name} in {duration:.2f}s")
                return result
            except Exception as e:
                duration = time.time() - start_time
                logger.error(f"Failed {name} after {duration:.2f}s: {str(e)}")
                raise
        
        # Return appropriate wrapper based on function type
        import asyncio
        if asyncio.iscoroutinefunction(func):
            return async_wrapper
        else:
            return sync_wrapper
            
    return decorator


# Context manager for logging sections
class LogSection:
    """Context manager for logging distinct sections of code"""
    
    def __init__(self, section_name: str, log_level: str = "INFO"):
        self.section_name = section_name
        self.log_level = log_level
        self.start_time = None
    
    def __enter__(self):
        import time
        self.start_time = time.time()
        logger.log(self.log_level, f"🚀 Starting {self.section_name}")
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        import time
        duration = time.time() - self.start_time
        
        if exc_type is None:
            logger.log(self.log_level, f"✅ Completed {self.section_name} in {duration:.2f}s")
        else:
            logger.error(f"❌ Failed {self.section_name} after {duration:.2f}s: {exc_val}")
        
        return False  # Don't suppress exceptions


# Utility functions for structured logging
def log_browser_info(browser_type: str, headless: bool, viewport: dict):
    """Log browser configuration information"""
    logger.info(f"Browser: {browser_type} (headless={headless}, viewport={viewport['width']}x{viewport['height']})")


def log_element_detection(elements: dict, method: str):
    """Log element detection results"""
    found_elements = list(elements.keys())
    logger.info(f"Element detection via {method}: found {found_elements}")


def log_captcha_attempt(service: str, captcha_type: str):
    """Log CAPTCHA solving attempt"""
    logger.info(f"Attempting to solve {captcha_type} CAPTCHA using {service}")


def log_login_attempt(url: str, username: str):
    """Log login attempt (without password)"""
    masked_username = username[:3] + "*" * (len(username) - 6) + username[-3:] if len(username) > 6 else "*" * len(username)
    logger.info(f"Login attempt for {masked_username} at {url}")


def log_error_with_context(error: Exception, context: dict):
    """Log error with additional context"""
    logger.error(f"Error: {str(error)}")
    for key, value in context.items():
        logger.error(f"  {key}: {value}")


# Configuration for different environments
def configure_for_production():
    """Configure logging for production environment"""
    logger.remove()
    
    # Only file logging in production
    logger.add(
        "logs/automation.log",
        format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {name}:{function}:{line} | {message}",
        level="INFO",
        rotation="50 MB",
        retention="30 days",
        compression="gz"
    )
    
    # Error file
    logger.add(
        "logs/errors.log",
        format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {name}:{function}:{line} | {message}",
        level="ERROR",
        rotation="10 MB",
        retention="30 days",
        compression="gz"
    )


def configure_for_development():
    """Configure logging for development environment"""
    logger.remove()
    
    # Console logging with colors
    logger.add(
        sys.stdout,
        format="<green>{time:HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}:{function}</cyan> | <level>{message}</level>",
        level="DEBUG",
        colorize=True
    )
    
    # File logging
    logger.add(
        "automation.log",
        format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {name}:{function}:{line} | {message}",
        level="DEBUG",
        rotation="10 MB",
        retention="3 days"
    )