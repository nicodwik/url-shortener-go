package user

import (
	"net/http"
	"url-shortener-go/helpers"
	UserRepository "url-shortener-go/repository/user"
	requests "url-shortener-go/requests/user"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	registerValidation := requests.RegisterValidation{Name: name, Email: email, Password: password}

	validationErrors := helpers.ValidateInput(registerValidation)
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError(validationErrors))
	}

	user, err := UserRepository.InsertNewUser(name, email, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("Email is used!"))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Data User", user))

}
