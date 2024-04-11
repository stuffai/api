package server

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/stuff-ai/api/internal/mongo"
	"github.com/stuff-ai/api/pkg/config"
	"github.com/stuff-ai/api/pkg/types/api"
)

var jwtKey = config.JWTKey()

// JWTClaims extends the standard jwt.Claims struct
type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken creates and signs a new JWT token for a given username
func generateToken(username string) (string, error) {
	claims := JWTClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// Set the expiration time
			ExpiresAt: time.Now().Add(24 * time.Hour * 365 * 50).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// login handles user login, returning a JWT token upon success
func login(c echo.Context) error {
	var req api.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Authenticate the user
	user, err := mongo.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	// Create and sign the token
	tokenString, err := generateToken(user.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not sign the token")
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

// signup handles the user signup process
func signup(c echo.Context) error {
	// Parse the request body to get signup details
	var requestBody api.SignupRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Minimum username length 2
	if len(requestBody.Username) < 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username must be at least 2 characters")
	}

	// Test if email is real
	if _, err := mail.ParseAddress(requestBody.Email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	// Pass the signup details to the mongo layer for user creation
	if _, err := mongo.InsertUser(c.Request().Context(), requestBody.Username, requestBody.Email, requestBody.Password); err != nil {
		// Log the error and return a generic error message to the client
		// You might want to handle different types of errors differently
		log.WithError(err).Error("api.signup: insertUser")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// Optionally, initiate the email verification process here

	// Return signed token.
	tokenString, err := generateToken(requestBody.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not sign the token")
	}

	// Return success response
	return c.JSON(http.StatusCreated, echo.Map{"token": tokenString})
}

// jwtMiddleware validates JWT tokens for protected routes
func jwtMiddleware(next echo.HandlerFunc, require bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		log.WithField("token", tokenString).Info("auth.jwt") // TODO(IMPORTANT): remove

		if err != nil || !token.Valid {
			if require {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			return next(c)
		}

		// add username to the context
		username := token.Claims.(*JWTClaims).Username
		uid, err := mongo.FindUserByName(c.Request().Context(), username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Backend failure")
		}
		c.Set("uid", uid)
		c.Set("username", username)

		// Token is valid, you can proceed with the request and also use the claims
		// For example, to get the username: claims.Username
		return next(c)
	}
}
