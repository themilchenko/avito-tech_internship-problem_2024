package httpAuth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/themilchenko/avito_internship-problem_2024/internal/config"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase

	cookieConfig config.CookieSettings
}

func NewAuthHandler(a domain.AuthUsecase, c config.CookieSettings) AuthHandler {
	return AuthHandler{
		authUsecase:  a,
		cookieConfig: c,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var recievedUser httpModels.User
	if err := c.Bind(&recievedUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	session, userID, err := h.authUsecase.SignUp(recievedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.SetCookie(h.makeHTTPCookie(session))

	return c.JSON(http.StatusOK, httpModels.ID{
		ID: userID,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var authUsr httpModels.User
	if err := c.Bind(&authUsr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	session, usrID, err := h.authUsecase.Login(authUsr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.SetCookie(h.makeHTTPCookie(session))

	return c.JSON(http.StatusOK, httpModels.ID{
		ID: usrID,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if err = h.authUsecase.Logout(cookie.Value); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	cookie.Expires = time.Now().AddDate(
		httpModels.DeleteExpire["year"],
		httpModels.DeleteExpire["month"],
		httpModels.DeleteExpire["day"],
	)
	c.SetCookie(h.makeHTTPCookie(cookie.Value))

	return c.JSON(http.StatusOK, httpModels.EmptyStruct{})
}

func (h *AuthHandler) Auth(c echo.Context) error {
	cookie, err := GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.authUsecase.Auth(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, httpModels.ID{
		ID: userID,
	})
}
