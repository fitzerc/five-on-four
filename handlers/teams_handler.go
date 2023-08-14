package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/fitzerc/five-on-four/utils"
	"github.com/labstack/echo/v4"
)

type TeamsHandler struct {
	TeamGuts             guts.TeamGuts
	TeamMessageBoardGuts guts.TeamMessageBoardGuts
	UserHandler          UserHandler
}

func (th TeamsHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/teams", th.AddTeam, th.UserHandler.MustBeAdmin())
	group.DELETE("/teams/:id", th.DeleteTeam, th.UserHandler.MustBeAdmin())
	group.GET("/teams/:id", th.GetTeamById, th.UserHandler.MustBeAdmin())
	group.GET("/teams", th.GetTeams, th.UserHandler.MustBeAdmin())
}

func (th TeamsHandler) AddTeam(c echo.Context) error {
	newTeam := new(data.Team)

	if err := c.Bind(newTeam); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingTeams, err := th.TeamGuts.GetByQuery("season_id = ? AND team_name = ?", newTeam.SeasonId, newTeam.TeamName)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	if len(existingTeams) > 0 && existingTeams[0].ID > 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "duplicate_team",
			ErrorDescription: "team already exists for this season",
		})
	}

	if err := th.TeamGuts.Add(*newTeam); err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	addedTeams, err := th.TeamGuts.GetByQuery("season_id = ? and team_name = ?", newTeam.SeasonId, newTeam.TeamName)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	teamMsgBoard := new(data.TeamMessageBoard)
	teamMsgBoard.TeamId = addedTeams[0].ID

	err = th.TeamMessageBoardGuts.Add(*teamMsgBoard)

	if err != nil {
		th.TeamGuts.Delete(utils.UintToString(addedTeams[0].ID))
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.String(http.StatusOK, "success")
}

func (th TeamsHandler) DeleteTeam(c echo.Context) error {
	id := c.Param("id")
	err := th.TeamGuts.Delete(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	teamMsgBoards, err := th.TeamMessageBoardGuts.GetByQuery("team_id = ?", id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	th.TeamMessageBoardGuts.Delete(utils.UintToString(teamMsgBoards[0].ID))

	return c.JSON(http.StatusOK, "success")
}

func (th TeamsHandler) GetTeamById(c echo.Context) error {
	id := c.Param("id")
	team, err := th.TeamGuts.GetById(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, team)
}

func (th TeamsHandler) GetTeams(c echo.Context) error {
	teams, err := th.TeamGuts.GetAll()

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, teams)
}
