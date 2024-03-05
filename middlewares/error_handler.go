package middlewares

import (
	"net/http"
	"url-shortener-go/helpers"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context, err error) error {
	return c.JSON(http.StatusUnauthorized, helpers.ResponseUnauthorized(err))
}
