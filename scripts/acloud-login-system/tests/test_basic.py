import pytest
from acloud_login_system import AILoginAutomator

@pytest.mark.asyncio
async def test_automator_initialization():
    """Test basic initialization of login automator"""
    automator = AILoginAutomator()
    assert automator is not None
    assert automator.config is not None
    assert hasattr(automator, 'perform_login')

def test_cli_help(runner):
    """Test CLI help command"""
    result = runner.invoke(cli.main, ['--help'])
    assert result.exit_code == 0
    assert 'AI-Powered Login Automation System' in result.output
