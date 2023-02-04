package mysql

import (
	"context"
)

// GetPublishedVideoIdsByUserId get list of published video ids
func GetPublishedVideoIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	res := make([]int64, 0)
	if err := DB.WithContext(ctx).Select("id").Model(&Video{}).Where("author_user_id = ?", userId).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
