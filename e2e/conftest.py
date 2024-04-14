import pytest

BASE_URL = "http://127.0.0.1:8081/api"


@pytest.fixture(scope="session", autouse=True)
def init():
    # yield
    pass


# @pytest.fixture(autouse=True)
# def destruct():
#     f = open('test.txt', 'a')
#     f.close()
