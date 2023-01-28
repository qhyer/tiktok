package db

import (
	"context"
)

// CreateVideo create video
func CreateVideo(ctx context.Context, video []*Video) error {
	return DB.WithContext(ctx).Create(video).Error
}

// GetPublishedVideosByUserId get list of published video
func GetPublishedVideosByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if err := DB.WithContext(ctx).Where(Video{AuthorUserId: userId}).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
