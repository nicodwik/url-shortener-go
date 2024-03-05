package user

import (
	"net/http"
	"url-shortener-go/helpers"
	UserRepository "url-shortener-go/repository/user"
	requests "url-shortener-go/requests/user"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	loginValidation := requests.LoginValidation{Email: email, Password: password}

	validationErrors := helpers.ValidateInput(loginValidation)
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError(validationErrors))
	}

	token, err := UserRepository.DoLogin(email, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("User not found!"))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Data User", token))

}
