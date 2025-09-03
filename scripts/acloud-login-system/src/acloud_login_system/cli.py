import click
from . import AILoginAutomator, get_logger

@click.command()
@click.option('--url', required=True, help='Login URL')
@click.option('--username', prompt=True, hide_input=False, help='Login username')
@click.option('--password', prompt=True, hide_input=True, help='Login password')
@click.option('--headless/--visible', default=True, help='Run browser in headless mode')
@click.version_option()
def main(url, username, password, headless):
    """AI-Powered Login Automation System"""
    logger = get_logger()
    
    try:
        automator = AILoginAutomator()
        logger.info(f"Starting login process for {url}")
        
        success = automator.perform_login(
            url=url,
            username=username,
            password=password,
            headless=headless
        )
        
        if success:
            logger.success("Login successful!")
            click.echo("Login successful!")
            return 0
        logger.error("Login failed")
        return 1
    except Exception as e:
        logger.critical(f"Critical error: {str(e)}")
        return 2

if __name__ == "__main__":
    main()
