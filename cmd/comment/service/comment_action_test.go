package service

import (
	"context"
	"reflect"
	"testing"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/comment"
	"tiktok/pkg/censor"
)

func TestCommentActionService_CommentAction(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *comment.DouyinCommentActionRequest
	}
	censor.Init()
	mysql.Init()
	text := "测试敏感词"
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *comment.Comment
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{req: &comment.DouyinCommentActionRequest{
				UserId:      1,
				VideoId:     1,
				ActionType:  1,
				CommentText: &text,
			}},
			wantErr: false,
		},
		{
			name:   "video not exist",
			fields: fields{context.Background()},
			args: args{req: &comment.DouyinCommentActionRequest{
				UserId:      1,
				VideoId:     0,
				ActionType:  1,
				CommentText: &text,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommentActionService{
				ctx: tt.fields.ctx,
			}
			got, err := s.CreateComment(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentActionService_DeleteCommentAction(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *comment.DouyinCommentActionRequest
	}
	commentIds := []int64{1, 0, 2}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{req: &comment.DouyinCommentActionRequest{
				UserId:     1,
				CommentId:  &commentIds[0],
				ActionType: 2,
			}},
			wantErr: false,
		},
		{
			name:   "comment not exist",
			fields: fields{context.Background()},
			args: args{req: &comment.DouyinCommentActionRequest{
				UserId:     1,
				CommentId:  &commentIds[1],
				ActionType: 2,
			}},
			wantErr: true,
		},
		{
			name:   "userId not match",
			fields: fields{context.Background()},
			args: args{req: &comment.DouyinCommentActionRequest{
				UserId:     2,
				CommentId:  &commentIds[2],
				ActionType: 2,
			}},
			wantErr: true,
		},
	}
	mysql.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommentActionService{
				ctx: tt.fields.ctx,
			}
			if err := s.DeleteComment(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
