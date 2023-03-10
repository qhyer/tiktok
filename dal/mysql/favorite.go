package mysql

import (
	"context"
	"time"

	"tiktok/pkg/constants"
	"tiktok/pkg/errno"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Favorite struct {
	gorm.Model
	Id        int64          `gorm:"column:id;primaryKey"`
	UserId    int64          `gorm:"column:user_id;uniqueIndex:uk_user_video_id"`
	VideoId   int64          `gorm:"column:video_id;uniqueIndex:uk_user_video_id"`
	CreatedAt time.Time      `gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (f *Favorite) TableName() string {
	return constants.FavoriteTableName
}

// GetFavoriteVideoIdsByUserId get favorite video ids by user id
func GetFavoriteVideoIdsByUserId(ctx context.Context, userId int64) ([]*Favorite, error) {
	res := make([]*Favorite, 0)
	if err := DB.WithContext(ctx).Model(&Favorite{}).Select("video_id", "created_at").Where("user_id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateFavorite user favorite video
func CreateFavorite(ctx context.Context, favorite *Favorite) (*Favorite, error) {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户新增喜欢
		res := tx.Model(&Favorite{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"deleted_at"}),
		}).Create(favorite)
		if res.Error != nil {
			return res.Error
		}

		// 创建失败
		if res.RowsAffected == 0 {
			return errno.FavoriteExistErr
		}

		// video 喜欢数+1
		res = tx.Model(&Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return errno.DatabaseOperationFailedErr
		}

		// 更新喜欢数失败
		if res.RowsAffected != 1 {
			return errno.VideoNotExistErr
		}

		return nil
	})
	return favorite, err
}

// DeleteFavorite user cancel favorite video
func DeleteFavorite(ctx context.Context, favorite *Favorite) (*Favorite, error) {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户取消喜欢
		res := tx.Where("video_id = ? and user_id = ?", favorite.VideoId, favorite.UserId).Delete(&favorite)
		if res.Error != nil {
			return errno.DatabaseOperationFailedErr
		}

		// 删除失败
		if res.RowsAffected != 1 {
			return errno.FavoriteNotExistErr
		}

		// video 喜欢数-1
		res = tx.Model(&Video{}).Where("id = ?", favorite.VideoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		// 更新喜欢数失败
		if res.RowsAffected != 1 {
			return errno.DatabaseOperationFailedErr
		}

		return nil
	})
	return favorite, err
}
