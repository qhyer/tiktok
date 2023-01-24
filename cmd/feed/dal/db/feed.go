package db

import (
	"context"
	"gorm.io/gorm"
	"tiktok/pkg/constants"
)

type Video struct {
	gorm.Model
}

func (v *Video) TableName() string {
	return constants.VideoTableName
}

// MGetVideos multiple get list of video
func MGetVideos(ctx context.Context, limit int, latestTime int64) ([]Video, error) {
	res := make([]*Video, 0)
	// TODO
	return res, nil
}
