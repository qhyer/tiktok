package db

import (
	"context"
)

// CreateVideo create video
func CreateVideo(ctx context.Context, video []*Video) error {
	return DB.WithContext(ctx).Create(video).Error
}

// GetPublishedVideoIdsByUserId get list of published video ids
func GetPublishedVideoIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	res := make([]int64, 0)
	if err := DB.WithContext(ctx).Select("id").Where(Video{AuthorUserId: userId}).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
