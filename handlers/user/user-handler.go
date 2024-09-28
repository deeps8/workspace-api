package userHandler

import (
	"net/http"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/labstack/echo/v4"
)

func FetchAllUsers(c echo.Context) error {
	user, ok := c.Get("user_details").(utils.UserDetails)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.Res{Message: "User not found in context", Ok: false})
	}

	users, err := database.GetAllUsers(user.Id)
	if err != nil {
		return c.JSON(http.StatusNotImplemented, err)
	}

	return c.JSON(http.StatusOK, utils.Res{Message: "Got ya all users", Ok: true, Data: users})
}
