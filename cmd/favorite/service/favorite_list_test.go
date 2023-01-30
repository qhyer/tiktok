package service

import (
	"context"
	"testing"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/rpc"
)

func TestFavoriteListService_FavoriteList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *favorite.DouyinFavoriteListRequest
	}
	mysql.Init()
	rpc.InitFeedRpc()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*feed.Video
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{context.Background()},
			args:    args{req: &favorite.DouyinFavoriteListRequest{UserId: 1}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteListService{
				ctx: tt.fields.ctx,
			}
			_, err := s.FavoriteList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
