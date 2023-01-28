package pack

import (
	"time"

	"tiktok/dal/db"
	"tiktok/kitex_gen/comment"
	"tiktok/kitex_gen/user"
)

// Comment pack comment
func Comment(com *db.Comment) *comment.Comment {
	if com == nil {
		return nil
	}
	return &comment.Comment{
		Id:         com.Id,
		Content:    com.Content,
		CreateDate: com.UpdatedAt.In(time.FixedZone("CST", 8*3600)).Format("01-02"),
		User:       &user.User{Id: com.UserId},
	}
}

// Comments pack list of comment
func Comments(cms []*db.Comment) []*comment.Comment {
	comments := make([]*comment.Comment, 0, len(cms))
	if len(cms) == 0 {
		return comments
	}

	// pack comment
	for _, c := range cms {
		if vp := Comment(c); vp != nil {
			comments = append(comments, vp)
		}
	}

	return comments
}
