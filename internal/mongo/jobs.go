package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/types"
)

func jobsCollection() *mongo.Collection {
	return db().Collection("jobs")
}

func InsertJob(ctx context.Context, promptID string) (string, error) {
	promptOID, err := primitive.ObjectIDFromHex(promptID)
	if err != nil {
		return "", err
	}

	result, err := jobsCollection().InsertOne(
		ctx,
		bson.D{{"promptID", promptOID}, {"state", 0}},
	)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func FindJobByID(ctx context.Context, jobID string) (*types.Job, error) {
	oid, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", oid}}
	job := new(types.Job)
	if err := jobsCollection().FindOne(ctx, filter).Decode(job); err != nil {
		return nil, err
	}
	return job, nil
}

func FindAllJobBuckets(ctx context.Context) ([]types.Bucket, error) {
	opts := options.Find().SetProjection(bson.D{{"bucket", 1}})
	cur, err := jobsCollection().Find(ctx, bson.D{{"state", 1}}, opts)
	if err != nil {
		return nil, err
	}
	var jobs []*types.Job
	if err = cur.All(ctx, &jobs); err != nil {
		return nil, err
	}
	buckets := make([]types.Bucket, len(jobs))
	for i, jobs := range jobs {
		buckets[i] = jobs.Bucket
	}
	return buckets, nil
}
