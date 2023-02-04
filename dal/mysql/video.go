package mysql

import (
	"context"
	"time"

	"tiktok/pkg/constants"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Id            int64          `gorm:"column:id;primaryKey" redis:"id"`
	AuthorUserId  int64          `gorm:"column:author_user_id;index:idx_author_user_id" redis:"author_user_id"`
	PlayUrl       string         `gorm:"column:play_url" redis:"play_url"`
	CoverUrl      string         `gorm:"column:cover_url" redis:"cover_url"`
	FavoriteCount int64          `gorm:"column:favorite_count;default:0" redis:"favorite_count"`
	CommentCount  int64          `gorm:"column:comment_count;default:0" redis:"comment_count"`
	Title         string         `gorm:"column:title" redis:"title"`
	CreatedAt     time.Time      `gorm:"column:created_at;index:idx_created_at" redis:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (v *Video) TableName() string {
	return constants.VideoTableName
}

// GetVideosByLatestTime get list of video
func GetVideosByLatestTime(ctx context.Context, limit int, latestTime int64) ([]*Video, error) {
	res := make([]*Video, 0, limit)
	if err := DB.WithContext(ctx).Limit(limit).Order("created_at desc").Find(&res, "created_at < ?", time.UnixMilli(latestTime)).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// MGetVideosByVideoIds get list of video by video ids
func MGetVideosByVideoIds(ctx context.Context, videoIds []int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if len(videoIds) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", videoIds).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
