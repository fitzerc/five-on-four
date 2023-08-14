package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type TeamMessageBoardsHandler struct {
	TeamMessageBoardGuts guts.TeamMessageBoardGuts
}

func (tmbh TeamMessageBoardsHandler) RegisterEndpoints(group *echo.Group) {
	group.GET("/teammessageboards/:id", tmbh.GetTeamById)
}

func (tmbh TeamMessageBoardsHandler) GetTeamById(c echo.Context) error {
	id := c.Param("id")
	teamMessageBoard, err := tmbh.TeamMessageBoardGuts.GetById(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, teamMessageBoard)
}
