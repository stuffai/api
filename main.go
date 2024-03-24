package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/stuff-ai/api/internal/bucket"
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
	e.GET("/feed", getFeed)
	e.POST("/prompts", postPrompts)
	e.GET("/prompts/rand", getPromptRand)
	e.POST("/jobs", postJobs)
	e.GET("/jobs/:id", getJobByID)
	e.GET("/jobs/:id/img", getJobImageURL)

	e.POST("/signup", signup)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func postJobs(c echo.Context) error {
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
	if err := rmq.Publish(context.Background(), []byte(jobID)); err != nil {
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

	return c.Redirect(http.StatusTemporaryRedirect, signedURL.String())
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
		img.URL = signedURL.String()
	}
	return c.JSON(http.StatusOK, feed)
}

// SignupRequestBody defines the expected structure of the signup request body
type SignupRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// signup handles the user signup process
func signup(c echo.Context) error {
	// Parse the request body to get signup details
	var requestBody SignupRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Pass the signup details to the mongo layer for user creation
	userID, err := mongo.InsertUser(c.Request().Context(), requestBody.Username, requestBody.Email, requestBody.Password)
	if err != nil {
		// Log the error and return a generic error message to the client
		// You might want to handle different types of errors differently
		log.WithError(err).Error("api.signup: insertUser")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// Optionally, initiate the email verification process here

	// Return success response
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfully",
		"user_id": userID,
	})
}
