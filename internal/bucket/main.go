package bucket

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

func init() {
	endpoint := "192.168.63.29:9000"
	accessKeyID := "9b2yTqkrIBlf2TPHDL24"
	secretAccessKey := "UX1fJraecnPD32W00mdpbFI5vi2MUzc6hn8lv7Jd"
	useSSL := false

	var err error
	if client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}); err != nil {
		panic("bucket.init: err: " + err.Error())
	}
}

func SignURL(ctx context.Context, bucket, key string) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", key))
	return client.PresignedGetObject(ctx, bucket, key, time.Hour, reqParams)
}
