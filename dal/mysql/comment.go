package mysql

import (
	"context"
	"time"

	"tiktok/pkg/constants"
	"tiktok/pkg/errno"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Id        int64          `gorm:"column:id;primaryKey"`
	UserId    int64          `gorm:"column:user_id" redis:"user_id"`
	VideoId   int64          `gorm:"column:video_id;index:idx_video_id"`
	Content   string         `gorm:"column:content" redis:"content"`
	CreatedAt time.Time      `gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c *Comment) TableName() string {
	return constants.CommentTableName
}

// GetCommentListByVideoId get list of video comment
func GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	if err := DB.WithContext(ctx).Where("video_id = ?", videoId).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// MGetCommentListByCommentId get list of video comment
func MGetCommentListByCommentId(ctx context.Context, commentIds []int64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	if err := DB.WithContext(ctx).Where("id in ?", commentIds).Order("created_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateComment user comment video
func CreateComment(ctx context.Context, comment *Comment) (*Comment, error) {
	return comment, DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 新增评论
		res := tx.Create(comment)
		if res.Error != nil {
			return res.Error
		}

		// 创建失败
		if res.RowsAffected == 0 {
			return errno.CommentExistErr
		}

		// video 评论数+1
		res = tx.Model(&Video{}).Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		// 更新评论数失败
		if res.RowsAffected != 1 {
			return errno.DatabaseOperationFailedErr
		}

		return nil
	})
}

// DeleteComment delete video comment action
func DeleteComment(ctx context.Context, comment *Comment) (*Comment, error) {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 找到要删除的评论
		delCom := tx.Where("id = ? and user_id = ?", comment.Id, comment.UserId).Take(&comment)
		if delCom.Error != nil {
			return errno.CommentNotExistErr
		}

		// 删除评论
		res := tx.Where("id = ?", comment.Id).Delete(&comment)
		if res.Error != nil {
			return errno.DatabaseOperationFailedErr
		}

		// 删除失败
		if res.RowsAffected != 1 {
			return errno.DatabaseOperationFailedErr
		}

		// video 评论数-1
		res = tx.Model(&Video{}).Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.Error != nil {
			return errno.DatabaseOperationFailedErr
		}

		// 更新评论数失败
		if res.RowsAffected != 1 {
			return errno.DatabaseOperationFailedErr
		}

		return nil
	})
	return comment, err
}
