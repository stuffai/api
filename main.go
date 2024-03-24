package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/rmq"
	"github.com/stuff-ai/api/pkg/types"
)

func main() {
	defer rmq.Shutdown()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/generate", generate)
	e.POST("/prompts", postPrompts)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func generate(c echo.Context) error {
	req := new(types.Prompt)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	b, err := json.Marshal(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := rmq.Publish(context.Background(), b); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK")
}

func postPrompts(c echo.Context) error {
	req := new(types.Prompt)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := mongo.AddPrompt(context.Background(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK")
}
