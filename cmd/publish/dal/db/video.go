package db

import (
	"context"
	"tiktok/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Id            int64          `gorm:"column:id"`
	AuthorUserId  int64          `gorm:"column:author_user_id"`
	PlayUrl       string         `gorm:"column:play_url"`
	CoverUrl      string         `gorm:"column:cover_url"`
	FavoriteCount int64          `gorm:"column:favorite_count;default:0"`
	CommentCount  int64          `gorm:"column:comment_count;default:0"`
	Title         string         `gorm:"column:title"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (v *Video) TableName() string {
	return constants.VideoTableName
}

// GetPublishedVideosByUserId get list of published video
func GetPublishedVideosByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if err := DB.WithContext(ctx).Where(Video{AuthorUserId: userId}).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
