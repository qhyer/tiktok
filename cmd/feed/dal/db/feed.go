package db

import (
	"context"
	"gorm.io/gorm"
	"tiktok/pkg/constants"
	"time"
)

type Video struct {
	gorm.Model
	Id            int64          `gorm:"column:id" json:"id"`
	AuthorUserId  int64          `gorm:"author_user_id"`
	PlayUrl       string         `gorm:"play_url"`
	CoverUrl      string         `gorm:"cover_url"`
	FavoriteCount int64          `gorm:"favorite_count"`
	CommentCount  int64          `gorm:"comment_count"`
	Title         string         `gorm:"title"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (v *Video) TableName() string {
	return constants.VideoTableName
}

// MGetVideos multiple get list of video
func MGetVideos(ctx context.Context, limit int, latestTime int64) ([]*Video, error) {
	res := make([]*Video, 0)
	// TODO
	if err := DB.WithContext(ctx).Limit(limit).Order("create_at desc").Find(&res, "create_at < ?", time.UnixMilli(latestTime)).Error; err != nil {
		return nil, err
	}
	return res, nil
}
