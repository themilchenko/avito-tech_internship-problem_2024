import requests
from requests.models import Response
from request_utils.banners_api import (
    Banner,
    BannerSearchRequest,
    create_banner_request,
    delete_banner_request,
    get_banner_request,
    get_banners_request,
)
from request_utils.auth_api import Auth, Token, login, logout
import pytest
import os

DEFAULT_ACTIVE_BANNER = Banner([1, 2], 23, "title", "text", "url", True)
DEFAULT_NOT_ACTIVE_BANNER = Banner([1, 2], 23, "title", "text", "url", False)


class BaseCase:
    token: Token

    @pytest.fixture(autouse=True, scope="function")
    def create_banner_enable(self, create_enable: bool, create_banner: Banner):
        if create_enable:
            resp = login(
                Auth(
                    str(os.getenv("ADMIN_USERNAME")),
                    str(os.getenv("ADMIN_PASSWORD")),
                    str(os.getenv("ADMIN_ROLE")),
                )
            )
            t = Token(str(resp.cookies.get("token")))

            get_response: requests.Response = get_banner_request(
                create_banner.banner_struct["tag_ids"][0],
                create_banner.banner_struct["feature_id"],
                False,
                t,
            )
            if get_response.status_code == 200:
                delete_banner_request(get_response.json()["banner_id"], t)

            create_banner_request(create_banner, t)
            logout(t)

            yield

            resp = login(
                Auth(
                    str(os.getenv("ADMIN_USERNAME")),
                    str(os.getenv("ADMIN_PASSWORD")),
                    str(os.getenv("ADMIN_ROLE")),
                )
            )
            t = Token(str(resp.cookies.get("token")))
            banners = get_banners_request(
                BannerSearchRequest(
                    create_banner.banner_struct["feature_id"],
                    create_banner.banner_struct["tag_ids"][0],
                ),
                t,
            )
            d = banners.json()[0]["banner_id"]
            delete_banner_request(d, t)
            logout(t)

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
