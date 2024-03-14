package routes

import (
	"url-shortener-go/config"
	RedirectionController "url-shortener-go/controllers/redirection"
	UserController "url-shortener-go/controllers/user"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {

	v1 := e.Group("/api/v1")
	redirection := v1.Group("/redirection")
	user := v1.Group("/user")

	{
		user.Use(echojwt.WithConfig(config.JwtConfig()))
		user.GET("", UserController.GetLoggedIn)
		user.GET("/redirections", UserController.GetRedirections)
		// user.GET("/:id", UserController.GetRedirections)
		user.GET("/:id/redirections", RedirectionController.GetAll)
	}

	redirection.GET("/:short_url", RedirectionController.Redirect)

	// redirection.GET("/", RedirectionController.Redirect)

	redirection.POST("/insert", RedirectionController.InsertNewUrl)
	redirection.POST("/batch-insert-seeder", RedirectionController.BatchInsertUrlsSeeder)

	v1.POST("/login", UserController.Login)
	v1.POST("/email-validator", UserController.EmailValidator)
	v1.POST("/register", UserController.Register)

}
