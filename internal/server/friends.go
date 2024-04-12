package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/queue"
	"github.com/stuff-ai/api/pkg/types"
	"github.com/stuff-ai/api/pkg/types/api"
)

func getFriends(c echo.Context) error {
	ctx := c.Request().Context()
	friends, err := mongo.FindFriends(ctx, c.Get("uid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := bucket.MaybeSignURLs(ctx, friends); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, friends)
}

func getFriendRequests(c echo.Context) error {
	ctx := c.Request().Context()
	reqs, err := mongo.FindFriendRequests(ctx, c.Get("uid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := bucket.MaybeSignURLs(ctx, reqs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, reqs)
}

func postFriendRequests(c echo.Context) error {
	ctx := c.Request().Context()
	uid := c.Get("uid")
	req := new(api.FriendRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	logger := log.WithFields(log.Fields{
		"method": "postFriendRequests",
		"user":   uid, "req": req,
	})
	// First check if the request hasn't already been made
	friendUID, err := mongo.FindUserByName(ctx, req.User)
	if err != nil {
		logger.WithError(err).Error("mongo.FindUserByName")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	exists, err := mongo.ExistsFriendRequest(ctx, friendUID, uid)
	if err != nil {
		logger.WithError(err).Error("mongo.ExistsFriendRequest")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	username := c.Get("username")
	if exists {
		// Check accepted status
		if req.Accepted {
			if err := mongo.AcceptFriendRequest(ctx, uid, friendUID); err != nil {
				logger.WithError(err).Error("mongo.AcceptFriendRequest")
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			notifID, err := mongo.InsertNotification(ctx, types.NotificationKindFriendAccepted, types.SignableMap{"id": uid, "user": username}, friendUID)
			if err != nil {
				logger.WithError(err).Error("mongo.insertNotification(friendAccept)")
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			if err := queue.PublishNotify(ctx, []byte(notifID.Hex())); err != nil {
				logger.WithError(err).Error("queue.PublishNotify(friendAccept)")
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			return c.JSON(http.StatusOK, "Accept OK")
		}
		// Reject friend request
		if err := mongo.RejectFriendRequest(ctx, friendUID, uid); err != nil {
			logger.WithError(err).Error("mongo.RejectFriendRequest")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, "Reject OK")
	}
	// Create friend request
	if err := mongo.InsertFriendRequest(ctx, uid, friendUID); err != nil {
		logger.WithError(err).Error("mongo.InsertFriendRequest")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	notifID, err := mongo.InsertNotification(ctx, types.NotificationKindFriendRequested, types.SignableMap{"id": uid, "user": username}, friendUID)
	if err != nil {
		logger.WithError(err).Error("mongo.insertNotification(friendRequested)")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := queue.PublishNotify(ctx, []byte(notifID.Hex())); err != nil {
		logger.WithError(err).Error("queue.PublishNotify(friendRequested)")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Request OK")
}
