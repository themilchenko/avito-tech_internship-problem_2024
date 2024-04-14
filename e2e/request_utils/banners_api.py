import requests
from request_utils.auth_api import Token
from conftest import BASE_URL


class Banner:
    def __init__(self, tag_ids, feature_id, title, text, url, is_active):
        self.banner_struct = {
            "tag_ids": tag_ids,
            "feature_id": feature_id,
            "content": {"title": title, "text": text, "url": url},
            "is_active": is_active,
        }


class BannerSearchRequest:
    def __init__(self, feature_id: int, tag_id: int, limit=0, offset=0):
        self.banner_params = {
            "feature_id": feature_id,
            "tag_id": tag_id,
            "limit": limit,
            "offset": offset,
        }


def get_banner_request(
    tag_id: int, feature_id: int, use_last_version: bool, token: Token
) -> requests.Response:
    params = {"tag_id": tag_id, "feature_id": feature_id, "use_last_version": False}
    return requests.get(BASE_URL + "/user_banner", params=params, cookies=token.token)


def create_banner_request(banner: Banner, token: Token):
    return requests.post(
        BASE_URL + "/banner", json=banner.banner_struct, cookies=token.token
    )


def delete_banner_request(banner_id: int, token: Token) -> requests.Response:
    return requests.delete(BASE_URL + f"/banner/{banner_id}", cookies=token.token)


def get_banners_request(
    banner_params: BannerSearchRequest, token: Token
) -> requests.Response:
    return requests.get(
        BASE_URL + "/banner", params=banner_params.banner_params, cookies=token.token
    )
