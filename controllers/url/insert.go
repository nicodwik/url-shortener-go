package url

import (
	"fmt"
	"net/http"
	"url-shortener-go/helpers"
	UrlRepository "url-shortener-go/repository/url"
	requests "url-shortener-go/requests/url"

	"github.com/labstack/echo/v4"
)

func InsertNewUrl(c echo.Context) error {
	shortUrl := c.FormValue("short_url")
	longUrl := c.FormValue("long_url")

	insertUrlValidation := requests.InsertUrlValidation{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}

	validationErrors := helpers.ValidateInput(insertUrlValidation)
	if len(validationErrors) > 0 {
		fmt.Println(validationErrors)

		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError(validationErrors))
	}

	findUrl, _ := UrlRepository.FindRedirection(shortUrl)
	if len(findUrl.ShortUrl) > 0 {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError([]string{"Short URL have been taken, please choose another!"}))
	}

	urlEntity, err := UrlRepository.InsertNewUrl(shortUrl, longUrl)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Url successfully created!", urlEntity))
}
