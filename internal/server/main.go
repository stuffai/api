package server

import (
	"net/http"

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

	// likes
	e.POST("/crafts/:id/likes", jwtMiddleware(postCraftLikes, true))
	e.DELETE("/crafts/:id/likes", jwtMiddleware(deleteCraftLikes, true))

	// comments
	e.GET("/crafts/:id/comments", jwtMiddleware(getCraftComments, false))
	e.POST("/crafts/:id/comments", jwtMiddleware(postCraftComments, true))

	// rank
	e.GET("/rank", getRank)
	e.POST("/rank", jwtMiddleware(postRank, true))

	// craft
	e.POST("/crafts", jwtMiddleware(postCrafts, true))

	// leaderboard
	e.GET("/leaderboard", getLeaderboard)

	// profile
	e.POST("/signup", signup)
	e.POST("/login", login)
	e.GET("/profile", jwtMiddleware(getProfile, true))
	e.PUT("/profile", jwtMiddleware(putProfile, true))
	e.POST("/profile/picture", jwtMiddleware(postProfilePicture, true))
	e.GET("/users/:name", jwtMiddleware(getUserProfile, false))
	e.POST("/fcm", jwtMiddleware(postFCM, true))
	e.DELETE("/fcm", jwtMiddleware(deleteFCM, true))
	e.GET("/notify", jwtMiddleware(notify, true))

	// friends
	e.GET("/friends", jwtMiddleware(getFriends, true))
	e.GET("/friends/requests", jwtMiddleware(getFriendRequests, true))
	e.POST("/friends/requests", jwtMiddleware(postFriendRequests, true))

	// notifs
	e.GET("/notifications", jwtMiddleware(getNotifications, true))
	// e.POST("/notifications", jwtMiddleware(postNotifications, true))
	e.PUT("/notifications/:id", jwtMiddleware(putNotification, true))

	// Private
	e.POST("/prompts", postPrompts)
	e.GET("/prompts/rand", getPromptRand)
	e.GET("/jobs/:id", getJobByID)
	e.GET("/jobs/:id/img", getJobImageURL)

	// HealthCheck
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })

	// Start server
	return e
}
