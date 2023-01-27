package db

import (
	"context"
	"gorm.io/gorm"
	"tiktok/pkg/constants"
	"time"
)

type Favorite struct {
	gorm.Model
	Id         int64          `gorm:"column:id"`
	UserId     int64          `gorm:"user_id"`
	VideoId    int64          `gorm:"video_id"`
	IsFavorite bool           `gorm:"is_favorite"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (f *Favorite) TableName() string {
	return constants.FavoriteTableName
}

// FavoriteList get list of user favorite videos
func FavoriteList(ctx context.Context, userId int64) ([]*Favorite, error) {
	res := make([]*Favorite, 0)
	if err := DB.WithContext(ctx).Where("user_id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func FavoriteAction(ctx context.Context, videoId int64, userId int64) {

}
