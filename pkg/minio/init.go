package minio

import (
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func Init() {
	// Initialize minio client object.
	c, err := minio.New(constants.OSSEndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(constants.OSSAccessKeyID, constants.OSSSecretAccessKey, ""),
	})
	if err != nil {
		klog.Errorf("minio client init failed %v", err)
	}
	Client = c
}
