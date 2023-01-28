package service

import (
	"context"
	"testing"

	"tiktok/dal/db"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/rpc"
)

func TestGetVideoService_GetVideosByVideoIdsAndCurrUserId(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest
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
			args: args{req: &feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest{
				VideoIds: []int64{1},
				UserId:   0,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GetVideoService{
				ctx: tt.fields.ctx,
			}
			_, err := s.GetVideosByVideoIdsAndCurrUserId(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideosByVideoIdsAndCurrUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
