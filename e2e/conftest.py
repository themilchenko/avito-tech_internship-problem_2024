import pytest
from dotenv import load_dotenv

BASE_URL = "https://milchenko.online/api/v2"


@pytest.fixture(scope="session", autouse=True)
def init():
    load_dotenv()

