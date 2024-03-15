package middlewares

import (
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func SuccessHandler(c echo.Context) {
	var loggedInUserId string

	token, ok := c.Get("token").(*jwtv4.Token)

	if ok && token.Valid {
		claims := token.Claims.(jwtv4.MapClaims)
		// loggedInUserId = claims.Id
		// fmt.Println(claims["id"])
		loggedInUserId = claims["id"].(string)
	}

	c.Set("loggedInUserId", loggedInUserId)
}
