package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"net/url"
	"os"
	"strconv"
	"time"

	"tiktok/dal/db"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/minio"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/satori/go.uuid"
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
	ruid := uuid.NewV4()
	fileName := strconv.FormatInt(time.Now().UnixMicro(), 62) + ruid.String()

	// 上传视频
	videoFileName := fileName + ".mp4"
	videoReader := bytes.NewReader(videoData)
	videoUploadInfo, err := minio.PutObject(s.ctx, constants.VideoBucketName, videoFileName, videoReader, int64(len(videoData)))
	if err != nil {
		klog.CtxErrorf(s.ctx, "upload file to oss failed %v", err)
		return errno.OSSUploadFailedErr
	}

	// 获取封面
	reqParams := make(url.Values)
	videoInfo, err := minio.PreSignedGetObject(s.ctx, constants.VideoBucketName, videoFileName, constants.OSSDefaultExpiry, reqParams)
	if err != nil {
		klog.CtxErrorf(s.ctx, "pre sign get object failed %v", err)
		return err
	}
	coverData, err := readFrameAsJpeg("http://" + constants.OSSEndPoint + videoInfo.RequestURI())
	if err != nil {
		return err
	}

	// 上传封面
	coverFileName := fileName + ".jpeg"
	coverReader := bytes.NewReader(coverData)
	coverUploadInfo, err := minio.PutObject(s.ctx, constants.CoverBucketName, coverFileName, coverReader, int64(len(coverData)))
	if err != nil {
		klog.CtxErrorf(s.ctx, "upload file to oss failed %v", err)
		return errno.OSSUploadFailedErr
	}

	// 在db插入结果
	err = db.CreateVideo(s.ctx, []*db.Video{
		{
			AuthorUserId: req.UserId,
			PlayUrl:      videoUploadInfo.Key,
			CoverUrl:     coverUploadInfo.Key,
			Title:        req.Title,
		},
	})
	if err != nil {
		klog.CtxFatalf(s.ctx, "db create video failed %v", err)
		return err
	}

	return nil
}

func readFrameAsJpeg(inFileName string) ([]byte, error) {
	reader := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}
