package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/stuff-ai/api/internal/bucket"
	"github.com/stuff-ai/api/internal/img"
	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/internal/queue"
	"github.com/stuff-ai/api/pkg/types"
)

func New() *echo.Echo {
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                        // Allows all origins
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, // Adjust methods as needed
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public
	e.POST("/signup", signup)
	e.POST("/login", login)
	e.GET("/feed", getFeed)
	e.POST("/crafts", jwtMiddleware(postCrafts))
	e.GET("/profile", jwtMiddleware(getProfile))
	e.PUT("/profile", jwtMiddleware(putProfile))
	e.POST("/profile/picture", jwtMiddleware(postProfilePicture))
	e.GET("/users/:name/feed", getUserFeed)

	// Private
	e.POST("/prompts", postPrompts)
	e.GET("/prompts/rand", getPromptRand)
	e.GET("/jobs/:id", getJobByID)
	e.GET("/jobs/:id/img", getJobImageURL)

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

func getProfile(c echo.Context) error {
	ctx := c.Request().Context()
	uid := c.Get("uid")

	// build profile (TODO optimize with aggregation or views or something)
	profile, err := mongo.GetUserProfile(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// sign profile picture
	if err := bucket.MaybeSignProfilePicture(ctx, profile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// craft count
	count, err := mongo.CountJobsForUser(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	profile.Crafts = int(count)
	// images
	imgs, err := mongo.FindImagesForUser(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := bucket.SignImages(ctx, imgs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	profile.Images = imgs

	return c.JSON(http.StatusOK, profile)
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
