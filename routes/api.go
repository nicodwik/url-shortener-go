package routes

import (
	"url-shortener-go/controllers/url"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {

	e.GET("/", url.GetAll)
	e.GET("/redirect/:short_url", url.Redirect)

	e.POST("/insert", url.InsertNewUrl)
	e.POST("/batch-insert-seeder", url.BatchInsertUrlsSeeder)

}
