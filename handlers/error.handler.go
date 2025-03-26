package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(code)
		}
	}

	// 🛡️ Важно: если 401 — добавляем WWW-Authenticate
	if code == http.StatusUnauthorized {
		c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	}

	c.Logger().Error(err)

	// Попробуем отдать кастомную страницу, если она есть
	errorPage := fmt.Sprintf("views/%d.html", code)
	if err := c.File(errorPage); err != nil {
		// Если страницы нет — просто отдадим текстовую ошибку
		c.String(code, fmt.Sprintf("%d %s", code, message))
	}
}
