package workspace

import (
	"net/http"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/labstack/echo/v4"
)

func CreateWorkspace(c echo.Context) error {
	user, ok := c.Get("user_details").(utils.UserDetails)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.Res{Message: "User not found in context", Ok: false})
	}

	details := utils.SpaceCreateDTO{}
	if err := c.Bind(&details); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{Message: err.Error(), Ok: false})
	}

	if len(details.Members) == 0 {
		return c.JSON(http.StatusBadRequest, utils.Res{Message: "Workspace should have atleast one member added", Ok: false})
	}

	details.Slug = utils.GenerateSlug(details.Name)
	details.Owner = user.Id

	err := database.InsertWorkspace(details)
	if err != nil {
		return c.JSON(http.StatusNotImplemented, utils.Res{Message: err.Error(), Ok: false})
	}
	// might need to return workspace details
	return c.JSON(http.StatusOK, utils.Res{Message: "Workspace Created", Ok: true})

}

func FetchAllWorkspaces(c echo.Context) error {

	spaces, err := database.GetAllWorkspace()
	if err != nil {
		return c.JSON(http.StatusNotImplemented, err)
	}
	return c.JSON(http.StatusOK, utils.Res{Message: "Got all Workspaces", Data: spaces, Ok: true})
}
