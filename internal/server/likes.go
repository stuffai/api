package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/mongo"
)

func postCraftLikes(c echo.Context) error {
	ctx := c.Request().Context()
	cid := c.Param("id")
	uid := c.Get("uid")

	_, err := mongo.FindImageByID(ctx, cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	if err := mongo.InsertLike(ctx, uid, cid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: handle notifications

	return c.JSON(http.StatusCreated, "Created")
}

func deleteCraftLikes(c echo.Context) error {
	ctx := c.Request().Context()
	cid := c.Param("id")
	uid := c.Get("uid")

	if err := mongo.DeleteLike(ctx, uid, cid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	_, err := mongo.FindImageByID(ctx, cid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: handle notifications

	return c.JSON(http.StatusOK, "OK")
}
