package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/queue"
	"github.com/stuff-ai/api/pkg/types"
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
	uid := c.Get("uid")
	comment := new(comment)
	if err := c.Bind(comment); err != nil || comment.Text == "" {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := mongo.InsertComment(ctx, uid, cid, comment.Text); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	craft, err := mongo.FindImageByID(ctx, cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	notifID, err := mongo.InsertNotification(ctx, types.NotificationKindCraftComment, types.SignableMap{
		"title":    craft.Title,
		"username": c.Get("username"),
		"comment":  comment.Text,
	}, craft.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := queue.PublishNotify(ctx, []byte(notifID.Hex())); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, "Created")
}
