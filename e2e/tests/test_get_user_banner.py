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

    @pytest.mark.parametrize(
        "user_auth, create_enable, create_banner, request_params, expected_response, expected_status_code, err",
        [
            (
                "admin",
                "true",
                Banner([1, 2], 1, "title", "text", "url", True),
                {
                    "tag_id": 1,
                    "feature_id": 1,
                    "use_last_version": False,
                },
                {
                    "title": "title",
                    "text": "text",
                    "url": "url",
                },
                200,
                "",
            ),
            (
                "user",
                "true",
                Banner([1, 2], 1, "title", "text", "url", True),
                {
                    "tag_id": 1,
                    "feature_id": 1,
                    "use_last_version": False,
                },
                {
                    "title": "title",
                    "text": "text",
                    "url": "url",
                },
                401,
                "this banner is not active right now",
            ),
            (
                "user",
                "true",
                Banner([1, 2], 1, "title", "text", "url", True),
                {
                    "tag_id": "a",
                    "feature_id": 1,
                    "use_last_version": False,
                },
                {
                    "title": "title",
                    "text": "text",
                    "url": "url",
                },
                400,
                'strconv.ParseUint: parsing "a": invalid syntax',
            ),
            (
                "user",
                "true",
                Banner([1, 2], 1, "title", "text", "url", True),
                {
                    "tag_id": 0,
                    "feature_id": 0,
                    "use_last_version": False,
                },
                {
                    "title": "title",
                    "text": "text",
                    "url": "url",
                },
                404,
                "failed to find item",
            ),
        ],
    )
    def test_get_active_user_banner_success(
        self,
        user_auth: str,
        create_enable: bool,
        create_banner: Banner,
        request_params,
        expected_response,
        expected_status_code,
        err,
    ):
        get_banner_response = get_banner_request(
            request_params["tag_id"],
            request_params["feature_id"],
            request_params["use_last_version"],
            self.token,
        )

        if err == "":
            assert get_banner_response.json() == expected_response
        else:
            assert get_banner_response.json()["message"] == err
