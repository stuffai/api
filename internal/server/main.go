package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/queue"
	"github.com/stuff-ai/api/pkg/types"
)

func New() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/feed", getFeed)
	e.POST("/prompts", postPrompts)
	e.GET("/prompts/rand", getPromptRand)
	e.POST("/crafts", jwtMiddleware(postCrafts))
	e.GET("/jobs/:id", getJobByID)
	e.GET("/jobs/:id/img", getJobImageURL)

	e.POST("/signup", signup)
	e.POST("/login", login)

	// Start server
	return e
}

// Handler
func postCrafts(c echo.Context) error {
	req := new(types.Prompt)

	// Parse req.
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Insert into DB.
	ctx := context.Background()
	promptID, err := mongo.InsertPrompt(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	jobID, err := mongo.InsertJob(ctx, promptID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Publish to queue.
	if err := queue.Publish(context.Background(), []byte(jobID)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{"jobID": jobID})
}

func postPrompts(c echo.Context) error {
	req := new(types.Prompt)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	id, err := mongo.InsertPrompt(context.Background(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func getPromptRand(c echo.Context) error {
	prompt, err := mongo.RandomPrompt(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prompt)
}

func getJobByID(c echo.Context) error {
	ctx := context.Background()
	prompt, err := mongo.FindJobByID(ctx, c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prompt)
}

func getJobImageURL(c echo.Context) error {
	ctx := context.Background()
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

func getFeed(c echo.Context) error {
	ctx := context.Background()
	feed, err := mongo.FindImages(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, img := range feed {
		signedURL, err := bucket.SignURL(ctx, img.Bucket.Name, img.Bucket.Key)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		img.URL = signedURL
	}
	return c.JSON(http.StatusOK, feed)
}
