package config

import (
	"strings"
	"time"
	"url-shortener-go/middlewares"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

type CustomClaim struct {
	jwt.RegisteredClaims
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GenerateJwtToken(id string, email string, name string) string {

	signedKeys := []byte("asdasd")

	claim := CustomClaim{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "asu",
		},
		id,
		email,
		name,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, _ := rawToken.SignedString(signedKeys)

	return token

}

func JwtConfig() echojwt.Config {
	return echojwt.Config{
		SigningKey:     []byte("asdasd"),
		ContextKey:     "token",
		SuccessHandler: middlewares.SuccessHandler,
		ErrorHandler:   middlewares.ErrorHandler,
	}
}

func GetUserIdByToken(c echo.Context) (string, error) {
	var token string
	var userId string

	authorization := c.Request().Header.Get("Authorization")
	if bearer := strings.Split(authorization, " "); len(bearer) == 2 {
		token = bearer[1]
	}

	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("asdasd"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := jwtToken.Claims.(*CustomClaim); ok {
		userId = claims.Id
	}

	return userId, nil
}
