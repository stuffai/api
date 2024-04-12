package server

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/pkg/types"
)

func notify(c echo.Context) error {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "stuffai-local-419503"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tokens, err := mongo.GetUserFCMTokens(ctx, c.Get("uid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.WithField("token", tokens).Info("server.notify")
	for _, entry := range tokens.Entries {
		msg := &messaging.Message{
			Data:  map[string]string{"title": "hello", "body": "notifications"},
			Token: entry.Token,
		}
		client, err := app.Messaging(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		if _, err := client.Send(ctx, msg); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, "OK")
}

func postFCM(c echo.Context) error {
	t := new(types.FCMToken)
	ctx := c.Request().Context()
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := mongo.UpdateUserFCMToken(ctx, c.Get("uid"), t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "OK")
}

func deleteFCM(c echo.Context) error {
	t := new(types.FCMToken)
	ctx := c.Request().Context()
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := mongo.DeleteUserFCMToken(ctx, c.Get("uid"), t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "OK")
}
