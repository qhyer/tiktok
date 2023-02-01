package service

import (
	"context"
	"testing"

	"tiktok/cmd/rpc"
	"tiktok/dal/mysql"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"
)

func TestPublishListService_PublishList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *publish.DouyinPublishListRequest
	}
	mysql.Init()
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
