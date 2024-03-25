package controllers

import (
	"net/http"
	"url-shortener-go/helpers"
	dashboard "url-shortener-go/repository"
	UserRepository "url-shortener-go/repository/user"

	"github.com/labstack/echo/v4"
)

type DashboardController struct {
	repo dashboard.DashboardContract
}

func NewDashboardController(dashboardRepo dashboard.DashboardContract) *DashboardController {
	return &DashboardController{dashboardRepo}
}

func (dc *DashboardController) Dashboard(c echo.Context) error {

	userId, _ := c.Get("loggedInUserId").(string)

	_, err := UserRepository.FindUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.ResponseNotFound("User not found!"))
	}

	dashboard, err := dc.repo.GetDashboardData(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.ResponseServerError("something went wrong!", err))
	}

	return c.JSON(http.StatusOK, helpers.ResponseOK("Data Dashboard", dashboard))
}
