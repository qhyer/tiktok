package service

import (
	"context"
	"testing"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/comment"
)

func TestCommentListService_CommentList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *comment.DouyinCommentListRequest
	}
	mysql.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*comment.Comment
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{ctx: context.Background()},
			args: args{&comment.DouyinCommentListRequest{
				VideoId: 1,
				UserId:  1,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommentListService{
				ctx: tt.fields.ctx,
			}
			_, err := s.CommentList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
