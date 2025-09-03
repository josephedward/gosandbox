#!/bin/bash

# AI Login Automation - Quick Setup Script
# This script automates the initial setup process

set -e

echo "🚀 AI Login Automation Setup Script"
echo "===================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if Python is installed
check_python() {
    print_step "Checking Python installation..."
    
    if ! command -v python3 &> /dev/null; then
        print_error "Python 3 is not installed. Please install Python 3.8 or higher."
        exit 1
    fi
    
    PYTHON_VERSION=$(python3 -c 'import sys; print(".".join(map(str, sys.version_info[:2])))')
    print_status "Found Python $PYTHON_VERSION"
    
    # Check if version is >= 3.8
    if [[ $(echo "$PYTHON_VERSION 3.8" | tr " " "\n" | sort -V | head -n1) != "3.8" ]]; then
        print_error "Python 3.8 or higher is required. Found: $PYTHON_VERSION"
        exit 1
    fi
}

# Check if pip is installed
check_pip() {
    print_step "Checking pip installation..."
    
    if ! command -v pip3 &> /dev/null; then
        print_error "pip3 is not installed. Please install pip."
        exit 1
    fi
    
    print_status "pip3 is available"
}

# Create virtual environment
create_venv() {
    print_step "Creating virtual environment..."
    
    if [ ! -d "venv" ]; then
        python3 -m venv venv
        print_status "Virtual environment created"
    else
        print_warning "Virtual environment already exists"
    fi
}

# Activate virtual environment
activate_venv() {
    print_step "Activating virtual environment..."
    
    source venv/bin/activate
    print_status "Virtual environment activated"
}

# Install requirements
install_requirements() {
    print_step "Installing Python requirements..."
    
    if [ ! -f "requirements.txt" ]; then
        print_error "requirements.txt not found!"
        exit 1
    fi
    
    pip install --upgrade pip
    pip install -r requirements.txt
    print_status "Python requirements installed"
}

# Install Playwright browsers
install_playwright() {
    print_step "Installing Playwright browsers..."
    
    playwright install chromium
    print_status "Playwright browsers installed"
}

# Create directory structure
create_directories() {
    print_step "Creating directory structure..."
    
    mkdir -p modules config tests logs screenshots
    
    # Create __init__.py files
    touch modules/__init__.py
    touch config/__init__.py
    touch tests/__init__.py
    
    print_status "Directory structure created"
}

# Setup environment file
setup_env_file() {
    print_step "Setting up environment file..."
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            print_status "Environment file created from template"
            print_warning "Please edit .env file with your credentials!"
        else
            # Create basic .env file
            cat > .env << EOF
# Authentication Credentials
LOGIN_USERNAME=your_email@example.com
LOGIN_PASSWORD=your_secure_password

# OpenAI Configuration
OPENAI_API_KEY=sk-your-openai-api-key-here
OPENAI_MODEL=gpt-4-turbo
OPENAI_MAX_TOKENS=1000

# CAPTCHA Service API Keys (at least one required)
CAPSOLVER_API_KEY=
CAPMONSTER_API_KEY=
TWOCAPTCHA_API_KEY=

# Browser Configuration
HEADLESS_MODE=false
BROWSER_TIMEOUT=30000
VIEWPORT_WIDTH=1920
VIEWPORT_HEIGHT=1080

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
STEALTH_MODE=true
RANDOM_DELAYS=true
EOF
            print_status "Basic environment file created"
            print_warning "Please edit .env file with your credentials!"
        fi
    else
        print_warning "Environment file already exists"
    fi
}

# Run basic tests
run_tests() {
    print_step "Running basic tests..."
    
    if [ -f "tests/test_automation.py" ]; then
        python -m pytest tests/test_automation.py::TestSettings::test_settings_defaults -v
        print_status "Basic tests passed"
    else
        print_warning "Test file not found, skipping tests"
    fi
}

# Check API connectivity
check_api_connectivity() {
    print_step "Checking API connectivity..."
    
    # Load environment variables
    if [ -f ".env" ]; then
        export $(cat .env | grep -v '#' | sed 's/\r$//' | awk '/=/ {print $0}')
    fi
    
    # Check OpenAI API
    if [ ! -z "$OPENAI_API_KEY" ] && [ "$OPENAI_API_KEY" != "sk-your-openai-api-key-here" ]; then
        python3 -c "
import openai
import os
try:
    client = openai.OpenAI(api_key=os.getenv('OPENAI_API_KEY'))
    # Simple API test
    response = client.models.list()
    print('✅ OpenAI API connection successful')
except Exception as e:
    print(f'❌ OpenAI API connection failed: {e}')
" || print_warning "OpenAI API check failed"
    else
        print_warning "OpenAI API key not configured"
    fi
}

# Main setup function
main() {
    echo ""
    print_step "Starting AI Login Automation setup..."
    echo ""
    
    # Run all setup steps
    check_python
    check_pip
    create_directories
    create_venv
    activate_venv
    install_requirements
    install_playwright
    setup_env_file
    run_tests
    check_api_connectivity
    
    echo ""
    print_status "Setup completed successfully! 🎉"
    echo ""
    echo "Next steps:"
    echo "1. Edit .env file with your credentials:"
    echo "   - Add your OpenAI API key"
    echo "   - Add your login credentials"
    echo "   - Add at least one CAPTCHA service API key"
    echo ""
    echo "2. Run the automation:"
    echo "   source venv/bin/activate  # Activate virtual environment"
    echo "   python main.py            # Run the automation"
    echo ""
    echo "3. For help and documentation:"
    echo "   cat README.md             # Read the documentation"
    echo "   python -m pytest tests/ -v  # Run all tests"
    echo ""
}

# Check if running with sudo (not recommended)
if [ "$EUID" -eq 0 ]; then
    print_warning "Running as root is not recommended for security reasons"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Run main setup
main