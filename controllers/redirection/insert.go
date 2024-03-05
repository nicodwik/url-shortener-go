package redirection

import (
	"fmt"
	"net/http"
	"url-shortener-go/config"
	"url-shortener-go/helpers"
	RedirectionRepository "url-shortener-go/repository/redirection"
	requests "url-shortener-go/requests/redirection"

	"github.com/labstack/echo/v4"
)

func InsertNewUrl(c echo.Context) error {
	shortUrl := c.FormValue("short_url")
	longUrl := c.FormValue("long_url")

	userId, _ := config.GetUserIdByToken(c)

	fmt.Println(userId)

	insertUrlValidation := requests.InsertRedirectionValidation{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}

	validationErrors := helpers.ValidateInput(insertUrlValidation)
	if len(validationErrors) > 0 {
		fmt.Println(validationErrors)

		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError(validationErrors))
	}

	findUrl, _ := RedirectionRepository.FindRedirection(shortUrl)
	if len(findUrl.ShortUrl) > 0 {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError([]string{"Short URL have been taken, please choose another!"}))
	}

	urlEntity, err := RedirectionRepository.InsertNewUrl(shortUrl, longUrl, userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Url successfully created!", urlEntity))
}
