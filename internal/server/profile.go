package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/img"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/pkg/types"
)

func getProfile(c echo.Context) error {
	uid := c.Get("uid")
	return _getProfile(c, uid, uid)
}

func putProfile(c echo.Context) error {
	profile := new(types.UserProfile)
	if err := c.Bind(profile); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := mongo.UpdateUserProfile(c.Request().Context(), c.Get("uid"), profile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func postProfilePicture(c echo.Context) error {
	ctx := c.Request().Context()

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	imgBuf, err := img.ProcessImage(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	bkt, key, err := bucket.UploadImage(ctx, c.Get("username").(string), imgBuf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := mongo.UpdateUserProfilePicture(ctx, c.Get("uid"), bkt, key); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func getUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	name := c.Param("name")

	// build profile (TODO optimize with aggregation or views or something)
	uid, err := mongo.FindUserByName(ctx, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return _getProfile(c, uid, c.Get("uid"))
}

func _getProfile(c echo.Context, getUID, requestUID interface{}) error {
	ctx := c.Request().Context()

	// build profile (TODO optimize with aggregation or views or something)
	profile, err := mongo.GetUserProfile(ctx, getUID, requestUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// sign profile picture
	if err := bucket.MaybeSignURL(ctx, profile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := bucket.SignURLs(ctx, profile.Images); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}
