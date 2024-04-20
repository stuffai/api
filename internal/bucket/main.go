package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/stuff-ai/api/pkg/config"
	"github.com/stuff-ai/api/pkg/types"
)

const (
	jpgContentType = "image/jpeg"
)

var (
	gcsClient   *storage.Client
	minioClient *minio.Client
)

func signURLGCS(bucket, object string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour),
	}

	u, err := gcsClient.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}
	return u, nil
}

func init() {
	ctx := context.Background()
	var err error
	if gcsClient, err = storage.NewClient(ctx); err != nil {
		panic(err)
	}
}

func SignURL(ctx context.Context, bucket, key string) (string, error) {
	return signURLGCS(bucket, key)
}

func SignURLs[T types.Signable](ctx context.Context, docs []T) error {
	for _, doc := range docs {
		bucket := doc.GetBucket()
		signedURL, err := SignURL(ctx, bucket.Name, bucket.Key)
		if err != nil {
			return err
		}
		doc.SetURL(signedURL)
	}
	return nil
}

func ppBucketKey(username string) (string, string) {
	return config.BucketName(), fmt.Sprintf("profiles/%s/%s", username, uuid.New().String())
}

// UploadImage uploads the image and returns the bucket and key where it's stored.
func UploadImage(ctx context.Context, username string, in *bytes.Buffer) (string, string, error) {
	bkt, key := ppBucketKey(username)

	return bkt, key, uploadImageGCS(ctx, in, bkt, key)
}

func uploadImageGCS(ctx context.Context, in *bytes.Buffer, bkt, key string) error {
	wc := gcsClient.Bucket(bkt).Object(key).NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, in); err != nil {
		return err
	}
	return nil
}

// MaybeSignURL inspects the bucket field of a Signable and signs url if exists
func MaybeSignURL(ctx context.Context, x types.Signable) error {
	bucket := x.GetBucket()
	if bucket.Key == "" {
		return nil
	}
	ppURL, err := SignURL(ctx, bucket.Name, bucket.Key)
	if err != nil {
		return err
	}
	x.SetURL(ppURL)
	return nil
}

// MaybeSignURLs inspects the ppBucket field of a UserProfile and signs the URL if it exists.
func MaybeSignURLs[T types.Signable](ctx context.Context, a []T) error {
	for _, x := range a {
		if err := MaybeSignURL(ctx, x); err != nil {
			return err
		}
	}
	return nil
}
