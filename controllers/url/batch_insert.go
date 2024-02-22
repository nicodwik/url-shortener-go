package url

import (
	"net/http"
	"strconv"
	"url-shortener-go/config"
	"url-shortener-go/helpers"
	UrlRepository "url-shortener-go/repository/url"

	"github.com/labstack/echo/v4"
)

func BatchInsertUrlsSeeder(c echo.Context) error {
	count := c.FormValue("count")

	if len(count) == 0 {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError([]string{"insert desired batch item count"}))
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError([]string{"insert number only"}))
	}

	urlEntities, err := config.RunUrlSeeder(countInt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ResponseServerError("Something went wrong!", err))
	}

	if err := UrlRepository.BatchInsertUrls(urlEntities); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ResponseServerError("Something went wrong!", err))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("URLs successfully created, total: "+count+" items", []string{}))
}
