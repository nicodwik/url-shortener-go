package user

import (
	"net/http"
	"url-shortener-go/helpers"
	RedirectionRepository "url-shortener-go/repository/redirection"
	UserRepository "url-shortener-go/repository/user"

	"github.com/labstack/echo/v4"
)

func GetRedirections(c echo.Context) error {
	userId, _ := c.Get("loggedInUserId").(string)

	user, err := UserRepository.FindUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("User not found!"))
	}

	redirections, _ := RedirectionRepository.GetAllByUserId(c, user.Id)

	return c.JSON(http.StatusOK, helpers.ResponseOK("Data Redirections", redirections))
}
