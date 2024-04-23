package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/queue"
)

func getCraftComments(c echo.Context) error {
	ctx := c.Request().Context()
	cid := c.Param("id")
	comments, err := mongo.FindComments(ctx, cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := bucket.MaybeSignURLs(ctx, comments); err != nil {
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
	// get notif listeners for craft and send notifs
	listeners, err := mongo.FindCraftListeners(ctx, cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	notifIDs, err := mongo.InsertNotificationsForCraftComment(
		ctx,
		craft,
		c.Get("username").(string),
		comment.Text,
		listeners,
		uid,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(notifIDs) > 0 {
		queue.PublishNotifyMany(ctx, notifIDs)
	}

	// insert commenting craft listener at end
	if err := mongo.MaybeInsertCraftListener(ctx, uid, cid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, "Created")
}
