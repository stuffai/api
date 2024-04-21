package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/mongo"
)

func getCraftComments(c echo.Context) error {
	ctx := c.Request().Context()
	cid := c.Param("id")
	comments, err := mongo.FindComments(ctx, cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, comments)
}

type comment struct {
	Text string `json:"text"`
}

func postCraftComments(c echo.Context) error {
	ctx := c.Request().Context()
	cid := c.Param("id")
	comment := new(comment)
	if err := c.Bind(comment); err != nil || comment.Text == "" {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := mongo.InsertComment(ctx, c.Get("uid"), cid, comment.Text); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, "Created")
}
