package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

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

func signURLMinio(bucket, key string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", key))
	url, err := minioClient.PresignedGetObject(context.Background(), bucket, key, time.Hour, reqParams)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func init() {
	ctx := context.Background()
	var err error
	if gcsClient, err = storage.NewClient(ctx); err != nil {
		panic(err)
	}

	// for local
	endpoint := "192.168.123.29:9000"
	accessKeyID := "9b2yTqkrIBlf2TPHDL24"
	secretAccessKey := "UX1fJraecnPD32W00mdpbFI5vi2MUzc6hn8lv7Jd"
	useSSL := false

	if minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}); err != nil {
		panic("bucket.init: err: " + err.Error())
	}
}

func SignURL(ctx context.Context, bucket, key string) (string, error) {
	if config.IsLocalEnv() {
		return signURLMinio(bucket, key)
	}
	return signURLGCS(bucket, key)
}

func SignImages(ctx context.Context, imgs []*types.Image) error {
	for _, img := range imgs {
		signedURL, err := SignURL(ctx, img.Bucket.Name, img.Bucket.Key)
		if err != nil {
			return err
		}
		img.URL = signedURL
	}
	return nil
}

func ppBucketKey(username string) (string, string) {
	return config.ProjectID(), fmt.Sprintf("profiles/%s/%s", username, uuid.New().String())
}

// UploadImage uploads the image and returns the bucket and key where it's stored.
func UploadImage(ctx context.Context, username string, in *bytes.Buffer) (string, string, error) {
	bkt, key := ppBucketKey(username)

	if config.IsLocalEnv() {
		return bkt, key, uploadImageMinio(ctx, in, bkt, key)
	}
	return bkt, key, uploadImageGCS(ctx, in, bkt, key)
}

func uploadImageMinio(ctx context.Context, in *bytes.Buffer, bkt, key string) error {
	_, err := minioClient.PutObject(ctx, bkt, key, in, int64(in.Len()), minio.PutObjectOptions{ContentType: jpgContentType})
	if err != nil {
		return err
	}
	return err
}

func uploadImageGCS(ctx context.Context, in *bytes.Buffer, bkt, key string) error {
	wc := gcsClient.Bucket(bkt).Object(key).NewWriter(ctx)
	if _, err := io.Copy(wc, in); err != nil {
		return err
	}
	return nil
}

// MaybeSignProfilePicture inspects the ppBucket field of a UserProfile and signs the URL if it exists.
func MaybeSignProfilePicture(ctx context.Context, pf types.Signable) error {
	bucket := pf.GetBucket()
	if bucket.Key == "" {
		return nil
	}
	ppURL, err := SignURL(ctx, bucket.Name, bucket.Key)
	if err != nil {
		return err
	}
	pf.SetURL(ppURL)
	return nil
}
