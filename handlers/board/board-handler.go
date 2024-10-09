package board

import (
	"log"
	"net/http"
	"work-space-backend/database"
	"work-space-backend/utils"

	"github.com/labstack/echo/v4"
)

func CreateBoard(c echo.Context) error {
	user, ok := c.Get("user_details").(utils.UserDetails)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.Res{Message: "User not found in context", Ok: false})
	}

	brd := utils.BoardDTO{}
	if err := c.Bind(&brd); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{Message: err.Error(), Ok: false})
	}

	brd.Slug = utils.GenerateSlug(brd.Name)
	brd.Owner = user.Id
	log.Printf("%+v", brd)
	res, err := database.InsertBoard(brd)
	if err != nil {
		return c.JSON(http.StatusNotImplemented, utils.Res{Message: err.Error(), Ok: false})
	}

	return c.JSON(http.StatusOK, utils.Res{Message: "Board created", Ok: true, Data: res})
}
