package service

import (
	"context"
	"reflect"
	"testing"
	"tiktok/cmd/feed/dal/db"
	"tiktok/cmd/feed/rpc"
	"tiktok/kitex_gen/feed"
	"time"
)

func TestFeedService_Feed(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *feed.DouyinFeedRequest
	}
	rpc.InitUserRpc()
	db.Init()
	ts := time.Now().UnixMilli()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*feed.Video
		want1   int64
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{req: &feed.DouyinFeedRequest{
				LatestTime: &ts,
				UserId:     0,
			}},
			want:    []*feed.Video{},
			want1:   1674726762000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FeedService{
				ctx: tt.fields.ctx,
			}
			got, got1, err := s.Feed(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Feed() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Feed() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
