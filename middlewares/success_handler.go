package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomClaim struct {
	jwt.RegisteredClaims
	Id string
}

func SuccessHandler(c echo.Context) {
	var loggedInUserId string
	token, ok := c.Get("token").(*jwt.Token)

	if ok && token.Valid {
		claims := token.Claims.(*CustomClaim)
		loggedInUserId = claims.Id
	}
	c.Set("loggedInUserId", loggedInUserId)
}
