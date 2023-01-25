package service

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PublishActionService struct {
	ctx context.Context
}

// NewPublishActionService new PublishActionService
func NewPublishActionService(ctx context.Context) *PublishActionService {
	return &PublishActionService{ctx: ctx}
}

// PublishAction publish video
func (s *PublishActionService) PublishAction(req *publish.DouyinPublishActionRequest) error {
	videoData := req.Data

	// 生成文件名
	ruid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	fileName := string(time.Now().UnixMicro()) + ruid.String()

	// Initialize minio client object.
	minioClient, err := minio.New(constants.OSSEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(constants.OSSAccessKeyID, constants.OSSSecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		klog.Errorf("minio client init failed %v", err)
	}

	// 上传视频
	videoFileName := fileName + ".mp4"
	reader := bytes.NewReader(videoData)
	videoInfo, err := minioClient.PutObject(s.ctx, constants.VideoBucketName, videoFileName, reader, int64(len(videoData)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		klog.Errorf("upload file to oss failed %v", err)
		return errno.OSSUploadFailedErr
	}

	// 获取封面
	reqParams := make(url.Values)
	videoFullUrl, err := minioClient.PresignedGetObject(s.ctx, constants.VideoBucketName, videoFileName, constants.OSSDefaultExpiry, reqParams)
	if err != nil {
		klog.Errorf("pre sign get object failed %v", err)
		return err
	}
	coverData, err := readFrameAsJpeg(videoFullUrl.RequestURI())
	klog.Infof("video url %v", videoFullUrl)
	if err != nil {
		return err
	}

	// 上传封面
	coverFileName := fileName + ".jpeg"
	coverInfo, err := minioClient.PutObject(s.ctx, constants.CoverBucketName, coverFileName, reader, int64(len(coverData)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		klog.Errorf("upload file to oss failed %v", err)
		return errno.OSSUploadFailedErr
	}

	// 在db插入结果

	return nil
}

func readFrameAsJpeg(inFileName string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
