# AI-Powered Login Automation System

An intelligent automation system that uses GPT-4 Vision and professional CAPTCHA solving services to automatically detect form elements and perform login operations on web pages.

## 🚀 Features

- **AI-Powered Element Detection**: Uses GPT-4 Vision to visually identify login form elements
- **Multi-Service CAPTCHA Solving**: Integrates with Capsolver, CapMonster, and 2Captcha
- **Robust Fallback System**: Multiple detection strategies ensure high success rates
- **Stealth Browser Automation**: Anti-detection measures for reliable operation
- **Comprehensive Logging**: Detailed logging for debugging and monitoring
- **Environment-Based Configuration**: Secure credential management
- **Extensive Testing**: Full test suite with mocks and integration tests

## 📋 Prerequisites

- Python 3.8 or higher
- Chrome/Chromium browser
- OpenAI API key (for GPT-4 Vision)
- At least one CAPTCHA solving service API key

## 🛠️ Installation

### 1. Clone/Download the Project

```bash
# Create project directory
mkdir ai-login-automation
cd ai-login-automation

# Download all the code files provided above into respective directories
```

### 2. Create Virtual Environment

```bash
# Create virtual environment
python -m venv venv

# Activate virtual environment
# On Windows:
venv\Scripts\activate
# On macOS/Linux:
source venv/bin/activate
```

### 3. Install Dependencies

```bash
# Install Python packages
pip install -r requirements.txt

# Install Playwright browsers
playwright install chromium
```

### 4. Configure Environment Variables

```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your credentials
nano .env  # or use your preferred editor
```

### Required Environment Variables

```env
# Authentication Credentials
LOGIN_USERNAME=your_email@example.com
LOGIN_PASSWORD=your_secure_password

# OpenAI Configuration
OPENAI_API_KEY=sk-your-openai-api-key-here

# CAPTCHA Service API Keys (at least one required)
CAPSOLVER_API_KEY=your_capsolver_api_key
CAPMONSTER_API_KEY=your_capmonster_api_key
TWOCAPTCHA_API_KEY=your_2captcha_api_key

# Optional Configuration
HEADLESS_MODE=false
LOG_LEVEL=INFO
```

## 🏃‍♂️ Usage

### Basic Usage

```bash
# Run the automation system
python main.py
```

### Programmatic Usage

```python
import asyncio
from config.settings import Settings
from main import AILoginAutomator

async def main():
    # Load configuration
    config = Settings()
    
    # Initialize automator
    automator = AILoginAutomator(config)
    
    # Perform login
    success = await automator.perform_login('https://example.com/login')
    
    if success:
        print("Login successful!")
    else:
        print("Login failed!")

# Run the automation
asyncio.run(main())
```

### Custom Configuration

```python
from config.settings import Settings

# Override specific settings
config = Settings()
config.HEADLESS_MODE = True
config.MAX_RETRIES = 5
config.RANDOM_DELAYS = True

automator = AILoginAutomator(config)
```

## 🧪 Testing

### Run All Tests

```bash
# Run full test suite
python -m pytest tests/ -v

# Run with coverage
python -m pytest tests/ --cov=. --cov-report=html

# Run specific test categories
python -m pytest tests/ -m "not integration" -v  # Skip integration tests
python -m pytest tests/ -m integration -v        # Only integration tests
```

### Test Individual Components

```bash
# Test settings
python -m pytest tests/test_automation.py::TestSettings -v

# Test AI detector
python -m pytest tests/test_automation.py::TestAIElementDetector -v

# Test CAPTCHA solver
python -m pytest tests/test_automation.py::TestCaptchaSolver -v
```

## 📁 Project Structure

```
ai-login-automation/
├── main.py                 # Main automation script
├── modules/
│   ├── __init__.py
│   ├── ai_detector.py      # GPT-4 Vision element detection
│   ├── captcha_solver.py   # Multi-service CAPTCHA solving
│   ├── element_mapper.py   # Fallback element detection
│   └── logger.py           # Centralized logging
├── config/
│   ├── __init__.py
│   └── settings.py         # Configuration management
├── tests/
│   ├── __init__.py
│   └── test_automation.py  # Comprehensive test suite
├── requirements.txt        # Python dependencies
├── .env.example           # Environment variables template
├── .env                   # Your environment variables (create this)
├── README.md              # This file
└── automation.log         # Log file (created automatically)
```

## 🔧 Configuration Options

### Browser Settings

```env
HEADLESS_MODE=false          # Run browser in headless mode
VIEWPORT_WIDTH=1920          # Browser window width
VIEWPORT_HEIGHT=1080         # Browser window height
BROWSER_TIMEOUT=30000        # Page load timeout (ms)
USER_AGENT=Mozilla/5.0...    # Custom user agent string
```

### Automation Settings

```env
MAX_RETRIES=3                # Maximum retry attempts
RETRY_DELAY=2                # Delay between retries (seconds)
ELEMENT_DETECTION_TIMEOUT=10 # Element detection timeout
SCREENSHOT_ON_FAILURE=true   # Take screenshot on failure
RANDOM_DELAYS=true           # Add random delays for human-like behavior
```

### Security Settings

```env
STEALTH_MODE=true            # Enable anti-detection measures
ENABLE_PROXY=false           # Use proxy server
PROXY_URL=http://user:pass@host:port  # Proxy configuration
```

### Logging Settings

```env
LOG_LEVEL=INFO               # Logging level (DEBUG, INFO, WARNING, ERROR)
LOG_FILE=automation.log      # Log file name
ENABLE_FILE_LOGGING=true     # Enable file logging
ENABLE_CONSOLE_LOGGING=true  # Enable console logging
```

## 🔐 API Keys and Services

### OpenAI API Key

1. Visit [OpenAI Platform](https://platform.openai.com/api-keys)
2. Create a new API key
3. Add it to your `.env` file as `OPENAI_API_KEY`

### CAPTCHA Solving Services

Choose one or more services for redundancy:

#### Capsolver (Recommended)
- Visit [Capsolver](https://www.capsolver.com/)
- High accuracy (99%+), fast solving
- Cost: ~$0.8/1000 reCAPTCHA v2 solves

#### CapMonster Cloud
- Visit [CapMonster](https://capmonster.cloud/)
- AI-only service, good performance
- Cost: ~$0.6/1000 tokens

#### 2Captcha
- Visit [2Captcha](https://2captcha.com/)
- Hybrid AI/human solving
- Cost: ~$1/1000 requests

## 📊 Monitoring and Debugging

### Log Files

The system generates detailed logs:

```bash
# View real-time logs
tail -f automation.log

# View error logs only
grep "ERROR" automation.log

# View specific session
grep "2024-01-01" automation.log
```

### Debug Mode

Enable debug mode for detailed information:

```env
LOG_LEVEL=DEBUG
HEADLESS_MODE=false
SCREENSHOT_ON_FAILURE=true
```

### Performance Monitoring

```python
from modules.logger import LogSection, log_performance

# Use context manager for timing
async with LogSection("Login Process"):
    result = await automator.perform_login(url)

# Use decorator for function timing
@log_performance("custom_function")
async def my_function():
    # Your code here
    pass
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Browser Launch Fails
```bash
# Reinstall browser engines
playwright install chromium

# Check if browser is available
playwright list
```

#### 2. Element Detection Fails
- Check if the page loads completely
- Verify the page structure hasn't changed
- Enable DEBUG logging to see detection attempts
- Try with `HEADLESS_MODE=false` to see what's happening

#### 3. CAPTCHA Solving Fails
- Verify API keys are correct
- Check service account balance
- Try different CAPTCHA services
- Ensure CAPTCHA is actually present on the page

#### 4. Login Success Detection Fails
- The system looks for common success indicators
- You may need to customize `_verify_login_success()` for specific sites
- Check if the login actually succeeded manually

### Getting Help

1. Enable debug logging: `LOG_LEVEL=DEBUG`
2. Check the log file for detailed error messages
3. Run tests to verify system components: `python -m pytest tests/ -v`
4. Take screenshots on failure: `SCREENSHOT_ON_FAILURE=true`

## ⚖️ Legal and Ethical Considerations

- **Respect Terms of Service**: Only use on websites where you have permission
- **Rate Limiting**: The system includes delays to avoid overwhelming servers
- **Legitimate Use Only**: Intended for automating your own account access
- **Security**: Never share API keys or credentials
- **Compliance**: Ensure usage complies with applicable laws and regulations

## 🔄 Advanced Usage

### Custom Element Detection

```python
# Extend the AIElementDetector for custom logic
class CustomDetector(AIElementDetector):
    async def detect_elements(self, page):
        # Custom detection logic here
        return await super().detect_elements(page)
```

### Multiple Site Support

```python
# Configure for different sites
sites = {
    'site1': {
        'url': 'https://site1.com/login',
        'username': 'user1@example.com',
        'password': 'password1'
    },
    'site2': {
        'url': 'https://site2.com/login', 
        'username': 'user2@example.com',
        'password': 'password2'
    }
}

for site_name, site_config in sites.items():
    # Update environment variables or settings
    config.LOGIN_USERNAME = site_config['username']
    config.LOGIN_PASSWORD = site_config['password']
    
    # Perform login
    success = await automator.perform_login(site_config['url'])
    print(f"{site_name}: {'✅' if success else '❌'}")
```

### Integration with Other Tools

```python
# Use with scheduling
import schedule
import time

def run_automation():
    asyncio.run(automator.perform_login(target_url))

# Schedule daily login
schedule.every().day.at("09:00").do(run_automation)

while True:
    schedule.run_pending()
    time.sleep(60)
```

## 📈 Performance Optimization

### Headless Mode
- Use `HEADLESS_MODE=true` for faster execution
- Only use visible mode for debugging

### Caching
- The system caches successful element detection patterns
- Subsequent runs on the same site are faster

### Parallel Processing
```python
# Run multiple automations concurrently
import asyncio

async def run_multiple_sites():
    tasks = []
    for url in urls:
        task = automator.perform_login(url)
        tasks.append(task)
    
    results = await asyncio.gather(*tasks)
    return results
```

## 🆕 Updates and Maintenance

### Keeping Dependencies Updated

```bash
# Update packages
pip install --upgrade -r requirements.txt

# Update browser
playwright install chromium
```

### Monitoring for Changes
- Web pages change frequently
- Monitor logs for detection failures
- Update selectors if needed
- Keep AI prompts current

## 📝 Contributing

Feel free to extend and improve the system:

1. Add new CAPTCHA services
2. Improve element detection algorithms
3. Add support for new authentication methods
4. Enhance error handling and recovery
5. Add new test cases

This system provides a robust foundation for AI-powered web automation while maintaining security and reliability standards.