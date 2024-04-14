import requests
from requests.models import Response
from request_utils.auth_api import Auth, Token, login, logout
import pytest
import os


class BaseCase:
    token: Token

    @pytest.fixture(autouse=True, scope="function")
    def setup(self, user_auth: str):
        resp: requests.Response = Response()
        if user_auth == "user":
            resp = login(
                Auth(
                    str(os.getenv("USER_USERNAME")),
                    str(os.getenv("USER_PASSWORD")),
                    str(os.getenv("USER_ROLE")),
                )
            )
        else:
            resp = login(
                Auth(
                    str(os.getenv("ADMIN_USERNAME")),
                    str(os.getenv("ADMIN_PASSWORD")),
                    str(os.getenv("ADMIN_ROLE")),
                )
            )

        self.token = Token(str(resp.cookies.get("token")))

        yield

        logout(self.token)
        self.token.token["token"] = ""
