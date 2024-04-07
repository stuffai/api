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

	// feed
	e.GET("/feed", getFeed)
	e.GET("/users/:name/feed", getUserFeed)

	// rank
	e.GET("/rank", getRank)
	e.POST("/rank", jwtMiddleware(postRank))

	// craft
	e.POST("/crafts", jwtMiddleware(postCrafts))

	// leaderboard
	e.GET("/leaderboard", getLeaderboard)

	// profile
	e.POST("/signup", signup)
	e.POST("/login", login)
	e.GET("/profile", jwtMiddleware(getProfile))
	e.PUT("/profile", jwtMiddleware(putProfile))
	e.POST("/profile/picture", jwtMiddleware(postProfilePicture))
	e.GET("/users/:name", optionalJwtMiddleware(getUserProfile))
	e.POST("/fcm", jwtMiddleware(postFCM))
	e.GET("/notify", jwtMiddleware(notify))

	// friends
	e.GET("/friends", jwtMiddleware(getFriends))
	e.GET("/friends/requests", jwtMiddleware(getFriendRequests))
	e.POST("/friends/requests", jwtMiddleware(postFriendRequests))

	// Private
	e.POST("/prompts", postPrompts)
	e.GET("/prompts/rand", getPromptRand)
	e.GET("/jobs/:id", getJobByID)
	e.GET("/jobs/:id/img", getJobImageURL)

	// Start server
	return e
}
