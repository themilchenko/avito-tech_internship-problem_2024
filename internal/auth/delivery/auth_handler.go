package httpAuth

import (
	"github.com/labstack/echo/v4"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(a domain.AuthUsecase) AuthHandler {
	return AuthHandler{
		authUsecase: a,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	return nil
}

func (h *AuthHandler) Login(c echo.Context) error {
	return nil
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return nil
}

func (h *AuthHandler) Auth(c echo.Context) error {
	return nil
}
