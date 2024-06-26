package httpAuth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	CookieName = "token"
)

func (h AuthHandler) makeHTTPCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:  CookieName,
		Value: sessionID,
		Expires: time.Now().
			AddDate(int(h.cookieConfig.ExpireDate.Years),
				int(h.cookieConfig.ExpireDate.Months),
				int(h.cookieConfig.ExpireDate.Days)),
		Secure:   false,
		HttpOnly: h.cookieConfig.HttpOnly,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
}

func GetCookie(c echo.Context) (*http.Cookie, error) {
	cookie, err := c.Cookie(CookieName)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}
