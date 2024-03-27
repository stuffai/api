package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stuff-ai/api/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func usersCollection() *mongo.Collection {
	return db().Collection("users")
}

func InsertUser(ctx context.Context, username, email, password string) (string, error) {
	// First, check if a user with the same email already exists
	var existingUser types.UserPrivate
	collection := usersCollection()
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{{"email": email}, {"username": username}}}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		if err == nil {
			// A user with this email already exists
			return "", errors.New("a user with this email already exists")
		}
		// Handle other potential errors from FindOne
		return "", err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new UserPrivate object
	newUser := types.UserPrivate{
		// ID is omitted to let MongoDB generate it
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		DTCreated:    time.Now(),
		DTModified:   time.Now(),
		Profile:      new(types.UserProfile),
	}

	// Insert the new user into the collection
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

// AuthenticateUser checks if a user exists with the given username and password
func AuthenticateUser(username, password string) (*types.UserPrivate, error) {
	collection := usersCollection() // Assume you have a function to get the users collection
	var user types.UserPrivate
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

type userID struct {
	ID primitive.ObjectID `bson:"_id"`
}

// FindUserByName returns the user ID of the given username as a Mongo ObjectID.
func FindUserByName(ctx context.Context, username string) (primitive.ObjectID, error) {
	var user userID
	if err := usersCollection().FindOne(ctx, bson.M{"username": username}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&user); err != nil {
		return primitive.NilObjectID, err
	}
	return user.ID, nil
}

// GetUserProfile
func GetUserProfile(ctx context.Context, uid interface{}) (*types.UserProfile, error) {
	var user types.UserPrivate
	if err := usersCollection().FindOne(ctx, bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"username": 1, "profile": 1})).Decode(&user); err != nil {
		log.WithError(err).Error("mongo.GetUserProfile")
		return nil, err
	}
	user.Profile.Username = user.Username // little hacky but we'll live
	return user.Profile, nil
}

// UpdateUserProfile updates a user's profile with the given input object
func UpdateUserProfile(ctx context.Context, uid interface{}, profile *types.UserProfile) error {
	_, err := usersCollection().UpdateByID(ctx, uid, bson.D{{"$set", bson.D{{"profile", profile}}}})
	return err
}

// UpdateUserProfilePicture updates the profile picture bucket information
func UpdateUserProfilePicture(ctx context.Context, uid interface{}, bkt, key string) error {
	_, err := usersCollection().UpdateByID(ctx, uid, bson.D{{"$set", bson.D{{"profile.ppBucket", types.Bucket{bkt, key}}}}})
	return err
}
