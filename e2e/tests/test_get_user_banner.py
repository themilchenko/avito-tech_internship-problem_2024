import json
import requests
import marshal
from tests.base_case import BaseCase
import pytest
from request_utils.banners_api import (
    BannerSearchRequest,
    create_banner_request,
    get_banner_request,
    delete_banner_request,
    Banner,
    get_banners_request,
)
from request_utils.auth_api import Token, Auth, signup, logout


class TestGetUserBanner(BaseCase):

    @pytest.mark.parametrize("user_auth", ["admin"])
    def test_get_active_user_banner_success(self, user_auth):
        get_response: requests.Response = get_banner_request(1, 23, False, self.token)
        if get_response.status_code != 404:
            delete_response = delete_banner_request(
                json.loads(get_response.content)["banner_id"], self.token
            )

        banner_to_create = Banner([1, 2], 23, "title", "text", "url", True)
        create_response: requests.Response = create_banner_request(
            banner_to_create, self.token
        )

        expected_response = {
            "title": "title",
            "text": "text",
            "url": "url",
        }
        get_banner_response = get_banner_request(
            banner_to_create.banner_struct["tag_ids"][1],
            banner_to_create.banner_struct["feature_id"],
            False,
            self.token,
        )
        assert get_banner_response.status_code == 200
        assert get_banner_response.json() == expected_response

        banners = get_banners_request(BannerSearchRequest(23, 1), self.token)
        d = banners.json()[0]['banner_id']
        delete_banner_request(d, self.token)

    # def test_get_not_active_
