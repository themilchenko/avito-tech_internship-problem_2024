package authMiddleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	httpAuth "github.com/themilchenko/avito_internship-problem_2024/internal/auth/delivery"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
)

type AuthMiddleware struct {
	authUsecase    domain.AuthUsecase
	bannersUsecase domain.BannersUsecase
}

func NewAuthMiddleware(a domain.AuthUsecase, b domain.BannersUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authUsecase:    a,
		bannersUsecase: b,
	}
}

func (m *AuthMiddleware) LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := httpAuth.GetCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrNoSession)
		}

		if _, err = m.authUsecase.Auth(cookie.Value); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrAuth)
		}

		return next(c)
	}
}

func (m *AuthMiddleware) AdminRequiured(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := httpAuth.GetCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrNoSession)
		}

		user, err := m.authUsecase.GetUserBySessionID(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrAuth)
		}
		if user.Role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, domain.ErrForbidden)
		}

		return next(c)
	}
}

func (m *AuthMiddleware) ActiveBannerRestriction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := httpAuth.GetCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrNoSession)
		}

		user, err := m.authUsecase.GetUserBySessionID(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, domain.ErrAuth)
		}

		tagID, err := strconv.ParseUint(c.QueryParam("tag_id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		featureID, err := strconv.ParseUint(c.QueryParam("feature_id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		usrBanner, err := m.bannersUsecase.GetUserBanner(tagID, featureID, false)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
}
