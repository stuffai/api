package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/pkg/types"
)

func getFeed(c echo.Context) error {
	ctx := c.Request().Context()
	feed, err := mongo.FindImages(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := _signImages(ctx, feed); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, feed)
}

func getUserFeed(c echo.Context) error {
	uname := c.Param("name")
	ctx := c.Request().Context()
	uid, err := mongo.FindUserByName(ctx, uname)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	feed, err := mongo.FindImagesForUser(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := _signImages(ctx, feed); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, feed)
}

func _signImages(ctx context.Context, feed []*types.Image) error {
	if err := bucket.SignImages(ctx, feed); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// TODO: This is likely to cause performance issues. Invest in a better solution.
	// - maybe a support field/endpoint containing a mapping of username/id to ppURLs
	for _, img := range feed {
		if err := bucket.MaybeSignProfilePicture(ctx, img.User); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
