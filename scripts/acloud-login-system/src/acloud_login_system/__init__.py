"""AI-powered login automation system for cloud credential management"""

__version__ = "0.1.0"

from .logger import setup_logger, get_logger, log_performance
from .ai_detector import AIElementDetector
from .captcha_solver import CaptchaSolver
from .main import AILoginAutomator

__all__ = [
    'setup_logger',
    'get_logger',
    'log_performance',
    'AIElementDetector',
    'CaptchaSolver',
    'AILoginAutomator'
]
