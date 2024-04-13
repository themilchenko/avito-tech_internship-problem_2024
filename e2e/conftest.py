import pytest
from utils.api_helpers import setup_api_client

@pytest.fixture
def api_client():
    return setup_api_client()
