package url

import (
	"log"
	"net/http"
	"url-shortener-go/helpers"
	UrlRepository "url-shortener-go/repository/url"

	"github.com/labstack/echo/v4"
)

func GetAll(c echo.Context) error {

	urls, err := UrlRepository.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("List Data", urls))
}
