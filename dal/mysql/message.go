package mysql

import (
	"context"
	"time"

	"tiktok/pkg/constants"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Id        int64          `gorm:"column:id;primaryKey" redis:"id"`
	UserId    int64          `gorm:"column:user_id;index:idx_user_touser_id" redis:"user_id"`
	ToUserId  int64          `gorm:"column:to_user_id;index:idx_user_touser_id" redis:"to_user_id"`
	Content   string         `gorm:"column:content" redis:"content"`
	CreatedAt time.Time      `gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c *Message) TableName() string {
	return constants.MessageTableName
}

// GetMessageListByUserId get list of message by two user ids
func GetMessageListByUserId(ctx context.Context, userId int64, toUserId int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := DB.WithContext(ctx).Where("user_id = ? and to_user_id = ?", userId, toUserId).Or("user_id = ? and to_user_id = ?", toUserId, userId).Order("created_at asc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// MGetMessageListByMessageId get list of message by message ids
func MGetMessageListByMessageId(ctx context.Context, msgIds []int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := DB.WithContext(ctx).Model(&Message{}).Where("id in ?", msgIds).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateMessage insert user message
func CreateMessage(ctx context.Context, messages []*Message) ([]*Message, error) {
	if err := DB.WithContext(ctx).Create(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
