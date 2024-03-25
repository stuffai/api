package bucket

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	endpoint := "192.168.63.29:9000"
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
	return signURLGCS(bucket, key)
}
