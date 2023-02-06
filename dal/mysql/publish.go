package mysql

import (
	"context"
)

// GetPublishedVideosByUserId get list of published video ids
func GetPublishedVideosByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if err := DB.WithContext(ctx).Model(&Video{}).Where("author_user_id = ?", userId).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
