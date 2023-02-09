package mysql

import (
	"context"
	"time"

	"tiktok/dal/neo4j"
	"tiktok/pkg/constants"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id            int64          `gorm:"column:id;primaryKey'" redis:"id"`
	UserName      string         `gorm:"column:username;unique_index:tk_username" redis:"name"`
	Password      string         `gorm:"column:password"`
	FollowCount   int64          `gorm:"-" redis:"follow_count"`
	FollowerCount int64          `gorm:"-" redis:"follower_count"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return constants.UserTableName
}

// CreateUser create user info
func CreateUser(ctx context.Context, users []*User) ([]*User, error) {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// mysql创建节点
		if err := tx.WithContext(ctx).Create(&users).Error; err != nil {
			return err
		}

		// 在neo4j中创建用户节点
		if err := neo4j.CreateUser(ctx, users); err != nil {
			return err
		}

		return nil
	})

	return users, err
}

// QueryUser query list of user info
func QueryUser(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
