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
	IsRead    bool           `gorm:"column:is_read;default:0"`
	CreatedAt time.Time      `gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c *Message) TableName() string {
	return constants.MessageTableName
}

// GetUnreadMessageListByUserId get list of video message and set message read
func GetUnreadMessageListByUserId(ctx context.Context, userId int64, toUserId int64) ([]*Message, error) {
	res := make([]*Message, 0)
	return res, DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 找到对方发给我的最早的未读消息
		var lstmsg *Message
		err := tx.Where("user_id = ? and to_user_id = ? and is_read = 0", toUserId, userId).Order("created_at asc").Take(&lstmsg).Error

		if err != nil {
			// 没找到记录返回空
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}

		// 找到从那条消息以来的所有消息
		err = tx.Where("user_id = ? and to_user_id = ? and created_at >= ?", userId, toUserId, lstmsg.CreatedAt).Or("user_id = ? and to_user_id = ? and created_at >= ?", toUserId, userId, lstmsg.CreatedAt).Order("created_at asc").Find(&res).Error
		if err != nil {
			return err
		}

		// 把这些消息标记为已读
		msgIds := make([]int64, 0, len(res))
		for _, m := range res {
			if m == nil {
				continue
			}
			// 消息是发给我的
			if m.ToUserId == userId {
				msgIds = append(msgIds, m.Id)
			}
		}

		// 更新为已读
		err = tx.Model(&Message{}).Where("id in ?", msgIds).Update("is_read", 1).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// CreateMessage insert user message
func CreateMessage(ctx context.Context, message []*Message) error {
	return DB.WithContext(ctx).Create(&message).Error
}
