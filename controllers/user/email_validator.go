package user

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"url-shortener-go/helpers"
	requests "url-shortener-go/requests/user"

	"github.com/labstack/echo/v4"
)

type EmailValidatorParam struct {
	ToEmail string `json:"to_email"`
}

type EmailValidatorResponse struct {
	IsReachable string `json:"is_reachable"`
}

func EmailValidator(c echo.Context) error {

	email := c.FormValue("email")

	UserEmailValidation := requests.UserEmailValidation{Email: email}
	validationErrors := helpers.ValidateInput(UserEmailValidation)
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, helpers.ResponseValidationError(validationErrors))
	}

	jsonParams, _ := json.Marshal(EmailValidatorParam{ToEmail: email})

	response, err := http.Post("https://emilia.aldion.dev/v0/check_email", "application/json", strings.NewReader(string(jsonParams)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.ResponseServerError("Something went wrong!", err))
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	emailValidatorResponse := EmailValidatorResponse{}
	_ = json.Unmarshal(body, &emailValidatorResponse)

	// CHECK EMAIL IS SAFE
	if emailValidatorResponse.IsReachable != "safe" {
		return c.JSON(http.StatusOK, helpers.ResponseValidationError([]string{"Email is invalid"}))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Email is safe", map[string]string{"email": email}))
}
