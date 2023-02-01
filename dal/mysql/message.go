package mysql

import (
	"context"
	"time"

	"tiktok/pkg/constants"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Id        int64          `gorm:"column:id;primaryKey"`
	UserId    int64          `gorm:"column:user_id;index:idx_user_touser_id"`
	ToUserId  int64          `gorm:"column:to_user_id;index:idx_user_touser_id"`
	Content   string         `gorm:"column:content"`
	CreatedAt time.Time      `gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c *Message) TableName() string {
	return constants.MessageTableName
}

// MessageList get list of video message
func MessageList(ctx context.Context, uid1 int64, uid2 int64, limit int) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := DB.WithContext(ctx).Where("user_id = ? and to_user_id = ?", uid1, uid2).Or("user_id = ? and to_user_id = ?", uid2, uid1).Order("created_at desc").Limit(limit).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateMessage insert user message
func CreateMessage(ctx context.Context, message []*Message) error {
	return DB.WithContext(ctx).Create(&message).Error
}
