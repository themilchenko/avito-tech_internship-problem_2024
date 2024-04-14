import requests

from conftest import BASE_URL


class Auth:
    def __init__(self, username: str, password: str, role: str):
        self.auth_user = {"username": username, "password": password, "role": role}


class Token:
    def __init__(self, token: str):
        self.token = {"token": token}


def signup(user: Auth):
    return requests.post(BASE_URL + "/signup", json=user.auth_user)


def login(user: Auth):
    return requests.post(BASE_URL + "/login", json=user.auth_user)


def logout(token: Token):
    return requests.delete(BASE_URL + "/logout", cookies=token.token)
