package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/rank"
	"github.com/stuff-ai/api/pkg/types"
	"github.com/stuff-ai/api/pkg/types/api"
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

func getRank(c echo.Context) error {
	ctx := c.Request().Context()
	imgs, err := mongo.FindImagesForRank(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := _signImages(ctx, imgs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, imgs)
}

func postRank(c echo.Context) error {
	ctx := c.Request().Context()
	r := new(api.RankRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ranks, err := mongo.FindImageRanksByIDs(ctx, r.Rank)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	// rank and coerce to map
	newRanks := rank.Rank(ranks)
	newRankMap := map[string]int{}
	for i, id := range r.Rank {
		newRankMap[id] = newRanks[i]
	}
	// submit to mongo
	if err := mongo.UpdateImageRanks(ctx, newRankMap); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return nil
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
