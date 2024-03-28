package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/users/:name", getUserProfile)
	e.GET("/users/:name/feed", getUserFeed)

	// Private
	e.POST("/prompts", postPrompts)
	e.GET("/prompts/rand", getPromptRand)
	e.GET("/jobs/:id", getJobByID)
	e.GET("/jobs/:id/img", getJobImageURL)

	// Start server
	return e
}
