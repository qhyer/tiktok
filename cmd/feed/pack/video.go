package pack

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
	"tiktok/cmd/feed/dal/db"
	"tiktok/cmd/feed/rpc"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
)

// Video pack video
func Video(video *db.Video, author *user.User, isFavorite bool) *feed.Video {
	if video == nil || author == nil {
		return nil
	}

	return &feed.Video{
		Id:            video.Id,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		Title:         video.Title,
		Author:        author,
		IsFavorite:    isFavorite,
	}
}

// Videos pack list of video
func Videos(ctx context.Context, vs []*db.Video, userId int64) ([]*feed.Video, int64, error) {
	var nextTime int64
	videos := make([]*feed.Video, 0, len(vs))

	if len(vs) == 0 {
		return videos, 0, nil
	}

	authorUserIds := make([]int64, 0, len(vs))

	for _, v := range vs {
		authorUserIds = append(authorUserIds, v.AuthorUserId)
	}

	// 查询作者的用户信息
	userInfoResponse, err := rpc.UserInfo(ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: authorUserIds,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "rpc query userinfo failed %v", userInfoResponse)
		return nil, 0, err
	}

	userInfoMap := make(map[int64]*user.User, 0)
	for _, us := range userInfoResponse.User {
		userInfoMap[us.Id] = &user.User{
			Id:            us.Id,
			Name:          us.Name,
			FollowCount:   us.FollowCount,
			FollowerCount: us.FollowerCount,
			IsFollow:      us.IsFollow,
		}
	}

	// TODO 查询用户是否给视频点赞
	userFavoriteMap := make(map[int64]bool, 0)

	// 找到最后一个视频的创建时间
	nextTime = vs[len(vs)-1].CreatedAt.UnixMilli()

	// pack video
	for _, v := range vs {
		if vp := Video(v, userInfoMap[v.AuthorUserId], userFavoriteMap[v.Id]); vp != nil {
			videos = append(videos, vp)
		}
	}
	// Initialize minio client object.
	minioClient, err := minio.New(constants.OSSEndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(constants.OSSAccessKeyID, constants.OSSSecretAccessKey, ""),
	})
	if err != nil {
		klog.Errorf("minio client init failed %v", err)
		return nil, nextTime, err
	}

	// oss sign
	for i, v := range videos {
		playUrl := v.PlayUrl
		coverUrl := v.CoverUrl
		reqParams := make(url.Values)
		videoInfo, err := minioClient.PresignedGetObject(ctx, constants.VideoBucketName, playUrl, constants.OSSDefaultExpiry, reqParams)
		if err != nil {
			klog.Errorf("pre sign get object failed %v", err)
			continue
		}
		coverInfo, err := minioClient.PresignedGetObject(ctx, constants.CoverBucketName, coverUrl, constants.OSSDefaultExpiry, reqParams)
		if err != nil {
			klog.Errorf("pre sign get object failed %v", err)
			continue
		}
		playUrl = constants.OSSBaseUrl + videoInfo.Path + "?" + videoInfo.RawQuery
		coverUrl = constants.OSSBaseUrl + coverInfo.Path + "?" + coverInfo.RawQuery
		videos[i].PlayUrl = playUrl
		videos[i].CoverUrl = coverUrl
	}

	return videos, nextTime, nil
}
