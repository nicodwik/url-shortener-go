package dashboard

import (
	"net/http"
	"url-shortener-go/helpers"
	DashboardRepository "url-shortener-go/repository"
	UserRepository "url-shortener-go/repository/user"

	"github.com/labstack/echo/v4"
)

func Dashboard(c echo.Context) error {

	userId, _ := c.Get("loggedInUserId").(string)

	user, err := UserRepository.FindUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("User not found!"))
	}

	dashboard, err := DashboardRepository.GetDashboardData(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.ResponseServerError("something went wrong!", err))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Data Dashboard", dashboard))
}
