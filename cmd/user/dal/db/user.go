package db

import (
	"context"
	"gorm.io/gorm"
	"tiktok/pkg/constants"
	"time"
)

type User struct {
	gorm.Model
	Id            int64          `gorm:"column:id"`
	UserName      string         `gorm:"column:username"`
	Password      string         `gorm:"column:password"`
	FollowCount   int64          `gorm:"column:follow_count;default:0"`
	FollowerCount int64          `gorm:"column:follower_count;default:0"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return constants.UserTableName
}

// MGetUsers multiple get list of user info
func MGetUsers(ctx context.Context, userIDs []int64) ([]*User, error) {
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUser create user info
func CreateUser(ctx context.Context, users []*User) error {
	return DB.WithContext(ctx).Create(users).Error
}

// QueryUser query list of user info
func QueryUser(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
