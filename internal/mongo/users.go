package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/stuff-ai/api/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func usersCollection() *mongo.Collection {
	return db().Collection("users")
}

func InsertUser(ctx context.Context, username, email, password string) (string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new UserPrivate object
	newUser := types.UserPrivate{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		DTCreated:    time.Now(),
		DTModified:   time.Now(),
		// ID is omitted to let MongoDB generate it
	}

	// Insert the new user into the collection
	collection := usersCollection()
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	// Assert that the returned ID is an ObjectId and convert it to string
	userID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("inserted ID is not an ObjectId")
	}

	return userID.Hex(), nil
}
