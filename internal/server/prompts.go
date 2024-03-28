package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/pkg/types"
)

func postPrompts(c echo.Context) error {
	req := new(types.Prompt)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	id, err := mongo.InsertPrompt(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func getPromptRand(c echo.Context) error {
	prompt, err := mongo.RandomPrompt(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prompt)
}
