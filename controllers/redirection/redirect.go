package redirection

import (
	"net/http"
	"url-shortener-go/helpers"
	RedirectionRepository "url-shortener-go/repository/redirection"

	"github.com/labstack/echo/v4"
)

type RedirectParams struct {
	ShortUrl string `param:"short_url"`
}

func Redirect(c echo.Context) error {

	var redirectParams RedirectParams

	if err := c.Bind(&redirectParams); err != nil {
		return err
	}

	url, err := RedirectionRepository.FindRedirection(redirectParams.ShortUrl)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("URL Not Found"))
	}

	if len(url.LongUrl) == 0 {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("URL Not Found"))
	}

	RedirectionRepository.IncementHitCount(url.ShortUrl)

	return c.JSON(http.StatusOK, helpers.ResponseOK("URL Found", map[string]string{"long_url": url.LongUrl}))
}
