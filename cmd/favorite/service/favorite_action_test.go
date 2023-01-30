package service

import (
	"context"
	"testing"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/favorite"
)

func TestFavoriteActionService_CancelFavoriteAction(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *favorite.DouyinFavoriteActionRequest
	}
	mysql.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{req: &favorite.DouyinFavoriteActionRequest{
				UserId:     1,
				VideoId:    1,
				ActionType: 2,
			}},
			wantErr: false,
		},
		{
			name:   "video not exist",
			fields: fields{context.Background()},
			args: args{req: &favorite.DouyinFavoriteActionRequest{
				UserId:     1,
				VideoId:    0,
				ActionType: 2,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteActionService{
				ctx: tt.fields.ctx,
			}
			if err := s.CancelFavoriteAction(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("CancelFavoriteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFavoriteActionService_FavoriteAction(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *favorite.DouyinFavoriteActionRequest
	}
	mysql.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{req: &favorite.DouyinFavoriteActionRequest{
				UserId:     1,
				VideoId:    1,
				ActionType: 1,
			}},
			wantErr: false,
		},
		{
			name:   "video not exist",
			fields: fields{context.Background()},
			args: args{req: &favorite.DouyinFavoriteActionRequest{
				UserId:     1,
				VideoId:    0,
				ActionType: 1,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FavoriteActionService{
				ctx: tt.fields.ctx,
			}
			if err := s.FavoriteAction(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("FavoriteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
