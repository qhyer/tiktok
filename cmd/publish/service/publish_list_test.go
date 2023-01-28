package service

import (
	"context"
	"testing"

	"tiktok/dal/db"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/rpc"
)

func TestPublishListService_PublishList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *publish.DouyinPublishListRequest
	}
	db.Init()
	rpc.InitUserRpc()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*feed.Video
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{req: &publish.DouyinPublishListRequest{
				UserId:   0,
				ToUserId: 1,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PublishListService{
				ctx: tt.fields.ctx,
			}
			_, err := s.PublishList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublishList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}