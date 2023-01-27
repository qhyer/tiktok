package db

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"time"
)

type Comment struct {
	gorm.Model
	Id        int64          `gorm:"column:id;primaryKey"`
	UserId    int64          `gorm:"column:user_id"`
	VideoId   int64          `gorm:"column:video_id;index:idx_video_id"`
	Content   string         `gorm:"column:content"`
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

func (c *Comment) TableName() string {
	return constants.CommentTableName
}

func (v *Video) TableName() string {
	return constants.VideoTableName
}

// CommentList get list of video comment
func CommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	if err := DB.WithContext(ctx).Where("video_id = ?", videoId).Order("updated_at desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CommentAction user comment video
func CommentAction(ctx context.Context, comment *Comment) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户新增评论
		res := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"deleted_at"}),
		}).Create(comment)
		if res.Error != nil {
			return res.Error
		}

		// 创建失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		// video 评论数+1
		res = tx.Model(&Video{}).Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		// 更新评论数失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		return nil
	})
}

// DeleteCommentAction delete video comment action
func DeleteCommentAction(ctx context.Context, comment *Comment) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户删除评论
		res := tx.Where("id = ? and user_id = ? and video_id = ?", comment.Id, comment.UserId, comment.VideoId).Delete(&comment)
		if res.Error != nil {
			return res.Error
		}

		// 删除失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		// video 评论数-1
		res = tx.Model(&Video{}).Where("id = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		// 更新评论数失败
		if res.RowsAffected != 1 {
			return errno.DBOperationFailedErr
		}

		return nil
	})
}
