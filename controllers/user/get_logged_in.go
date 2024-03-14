package user

import (
	"net/http"
	"url-shortener-go/helpers"
	UserRepository "url-shortener-go/repository/user"

	"github.com/labstack/echo/v4"
)

func GetLoggedIn(c echo.Context) error {
	userId, _ := c.Get("loggedInUserId").(string)

	user, err := UserRepository.FindUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("User not found!"))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Data User", user))
}
