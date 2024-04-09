package authMiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpAuth "github.com/themilchenko/avito_internship-problem_2024/internal/auth/delivery"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
)

type AuthMiddleware struct {
	authUsecase domain.AuthUsecase
}

func NewAuthMiddleware(a domain.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authUsecase: a,
	}
}

func (m *AuthMiddleware) LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := httpAuth.GetCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		if _, err = m.authUsecase.Auth(cookie.Value); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		return next(c)
	}
}

func (m *AuthMiddleware) AdminRequiured(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := httpAuth.GetCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		user, err := m.authUsecase.GetUserBySessionID(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		if user.Role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, err)
		}

		return next(c)
	}
}
