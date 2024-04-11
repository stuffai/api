package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
)

func getNotifications(c echo.Context) error {
	uid := c.Get("uid")
	ctx := c.Request().Context()
	notifs, err := mongo.GetNotifications(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := bucket.MaybeSignURLs(ctx, notifs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, notifs)
}

/*
func postNotifications(c echo.Context) error {
	notif := new(types.Notification)
	if err := c.Bind(notif); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	uid := c.Get("uid")
	if err := mongo.InsertNotification(c.Request().Context(), notif.Kind, notif.Data, uid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "OK")
}
*/

func putNotification(c echo.Context) error {
	notifID := c.Param("id")
	if err := mongo.UpdateNotificationRead(c.Request().Context(), notifID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "OK")
}
