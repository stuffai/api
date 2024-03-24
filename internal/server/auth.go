package server

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stuff-ai/api/internal/mongo" // Use the correct path
)

var jwtKey = []byte("your_secret_key") // This should be a secret key

// JWTClaims extends the standard jwt.Claims struct
type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// login handles user login, returning a JWT token upon success
func login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Authenticate the user
	user, err := mongo.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not sign the token")
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
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
