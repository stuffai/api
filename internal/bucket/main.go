package bucket

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var svc *s3.S3

func init() {
	endpoint := os.Getenv("GCS_URL")
	accessKeyID := os.Getenv("GCS_ACCESS_KEY")
	secretAccessKey := os.Getenv("GCS_SECRET_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		panic(err)
	}
	svc = s3.New(sess)
}

func SignURL(ctx context.Context, bucket, key string) (string, error) {
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	signedURL, err := req.Presign(24 * time.Hour)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}
