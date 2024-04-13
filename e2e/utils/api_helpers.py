import requests

def get_banner(tag_id, feature_id, token):
    url = "http://example.com/user_banner"
    params = {
        "tag_id": tag_id,
        "feature_id": feature_id,
        "use_last_revision": False
    }
    headers = {
        "token": token
    }
    response = requests.get(url, params=params, headers=headers)
    return response.json()
