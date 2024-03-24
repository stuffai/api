package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/stuff-ai/api/internal/rmq"
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

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

type Request struct {
	Title string `json:"title"`
	Prompt string `json:"prompt"`
}

// Handler
func generate(c echo.Context) error {
	req := new(Request)
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
