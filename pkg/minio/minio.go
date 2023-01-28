package minio

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

func PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64) (info minio.UploadInfo, err error) {
	return client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
}

func PreSignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	return client.PresignedGetObject(ctx, bucketName, objectName, expires, reqParams)
}
