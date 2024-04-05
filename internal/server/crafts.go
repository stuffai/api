package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/queue"
	"github.com/stuff-ai/api/pkg/types"
)

// Handler
func postCrafts(c echo.Context) error {
	req := new(types.Prompt)

	// Parse req.
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Insert into DB.
	ctx := c.Request().Context()
	promptID, err := mongo.InsertPrompt(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	jobID, err := mongo.InsertJob(ctx, c.Get("uid"), promptID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Publish to queue.
	if err := queue.PublishGenerate(context.Background(), []byte(jobID)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{"jobID": jobID})
}

func getJobByID(c echo.Context) error {
	prompt, err := mongo.FindJobByID(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prompt)
}

func getJobImageURL(c echo.Context) error {
	ctx := c.Request().Context()
	prompt, err := mongo.FindJobByID(ctx, c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	signedURL, err := bucket.SignURL(ctx, prompt.Bucket.Name, prompt.Bucket.Key)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, signedURL)
}
