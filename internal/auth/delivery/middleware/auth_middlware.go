package authMiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
)

type AuthMiddleware struct {
	authRepository domain.AuthUsecase
}

func NewAuthMiddleware(a domain.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authRepository: a,
	}
}

func (m *AuthMiddleware) LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
