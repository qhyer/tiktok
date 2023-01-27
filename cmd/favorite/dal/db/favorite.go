package db

import (
	"context"
	"gorm.io/gorm/clause"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	Id        int64          `gorm:"column:id;primaryKey"`
	UserId    int64          `gorm:"column:user_id;uniqueIndex:uk_user_video_id"`
	VideoId   int64          `gorm:"column:video_id;uniqueIndex:uk_user_video_id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index:idx_updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

type Video struct {
	gorm.Model
	Id            int64          `gorm:"column:id;primaryKey"`
	AuthorUserId  int64          `gorm:"column:author_user_id;index:idx_author_user_id"`
	PlayUrl       string         `gorm:"column:play_url"`
	CoverUrl      string         `gorm:"column:cover_url"`
	FavoriteCount int64          `gorm:"column:favorite_count;default:0"`
	CommentCount  int64          `gorm:"column:comment_count;default:0"`
	Title         string         `gorm:"column:title"`
	CreatedAt     time.Time      `gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (f *Favorite) TableName() string {
	return constants.FavoriteTableName
}

func (v *Video) TableName() string {
	return constants.VideoTableName
}

// FavoriteList get list of user favorite videos
func FavoriteList(ctx context.Context, userId int64) ([]*Favorite, error) {
	res := make([]*Favorite, 0)
	if err := DB.WithContext(ctx).Where("user_id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// FavoriteAction user favorite video
func FavoriteAction(ctx context.Context, favorite *Favorite) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户新增喜欢
		res := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"deleted_at"}),
		}).Create(favorite)
		if res.Error != nil {
			return res.Error
		}

		// 创建失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		// video 喜欢数+1
		res = tx.Model(&Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		// 更新喜欢数失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		return nil
	})
}

// CancelFavoriteAction user cancel favorite video
func CancelFavoriteAction(ctx context.Context, favorite *Favorite) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户取消喜欢
		res := tx.Where("video_id = ? and user_id = ?", favorite.VideoId, favorite.UserId).Delete(&favorite)
		if res.Error != nil {
			return res.Error
		}

		// 删除失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		// video 喜欢数-1
		res = tx.Model(&Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		// 更新喜欢数失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		return nil
	})
}
