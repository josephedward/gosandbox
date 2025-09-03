# AI Login Automation System - Complete Code Package

## Project Structure
```
ai-login-automation/
├── main.py                 # Main automation script
├── modules/
│   ├── __init__.py
│   ├── ai_detector.py      # AI-based element detection
│   ├── captcha_solver.py   # CAPTCHA solving logic
│   ├── element_mapper.py   # Coordinate to selector mapping
│   └── logger.py           # Logging utilities
├── config/
│   ├── __init__.py
│   └── settings.py         # Configuration management
├── tests/
│   ├── __init__.py
│   └── test_automation.py  # Unit tests
├── requirements.txt        # Python dependencies
├── .env.example           # Environment variables template
├── README.md              # Setup and usage instructions
└── docker/
    ├── Dockerfile         # Container setup
    └── docker-compose.yml # Multi-service deployment
```

## Installation & Setup

### Prerequisites
- Python 3.8+ 
- Chrome/Chromium browser
- OpenAI API key
- Capsolver API key (or alternative CAPTCHA service)

### Quick Start
```bash
# Clone or create project directory
mkdir ai-login-automation && cd ai-login-automation

# Create virtual environment
python -m venv venv
source venv/bin/activate  # Linux/Mac
# venv\Scripts\activate    # Windows

# Install dependencies
pip install -r requirements.txt

# Install browser engines
playwright install chromium

# Copy environment template and configure
cp .env.example .env
# Edit .env with your credentials
```

## Main Configuration Files

### requirements.txt
```
playwright==1.45.0
openai==1.12.0
capsolver==1.0.10
python-dotenv==1.0.1
opencv-python==4.9.0.80
pillow==10.2.0
requests==2.31.0
beautifulsoup4==4.12.3
lxml==5.1.0
pytest==8.0.0
pytest-asyncio==0.23.5
loguru==0.7.2
tenacity==8.2.3
pydantic==2.6.1
python-multipart==0.0.9
aiofiles==23.2.1
httpx==0.26.0
```

### .env.example
```env
# Authentication Credentials
LOGIN_USERNAME=your_email@example.com
LOGIN_PASSWORD=your_secure_password

# OpenAI Configuration
OPENAI_API_KEY=sk-your-openai-api-key-here
OPENAI_MODEL=gpt-4-turbo
OPENAI_MAX_TOKENS=1000

# CAPTCHA Service API Keys (choose one or multiple for redundancy)
CAPSOLVER_API_KEY=your_capsolver_api_key
CAPMONSTER_API_KEY=your_capmonster_api_key
TWOCAPTCHA_API_KEY=your_2captcha_api_key

# Browser Configuration
HEADLESS_MODE=false
BROWSER_TIMEOUT=30000
VIEWPORT_WIDTH=1920
VIEWPORT_HEIGHT=1080
USER_AGENT=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36

# Automation Settings
MAX_RETRIES=3
RETRY_DELAY=2
SCREENSHOT_ON_FAILURE=true
ELEMENT_DETECTION_TIMEOUT=10

# Logging Configuration
LOG_LEVEL=INFO
LOG_FILE=automation.log
ENABLE_FILE_LOGGING=true
ENABLE_CONSOLE_LOGGING=true

# Security Settings
ENABLE_PROXY=false
PROXY_URL=http://username:password@proxy-server:port
STEALTH_MODE=true
RANDOM_DELAYS=true
```