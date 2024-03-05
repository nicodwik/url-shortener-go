package redirection

import (
	"log"
	"net/http"
	"url-shortener-go/helpers"
	RedirectionRepository "url-shortener-go/repository/redirection"

	"github.com/labstack/echo/v4"
)

func GetAll(c echo.Context) error {

	urls, err := RedirectionRepository.GetAll(c)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("List Data", urls))
}
