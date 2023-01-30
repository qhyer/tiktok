package service

import (
	"context"
	"log"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
)

func TestFollowListService_FollowList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *relation.DouyinRelationFollowListRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*user.User
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{ctx: context.Background()},
			args:    args{req: &relation.DouyinRelationFollowListRequest{UserId: 1, ToUserId: 1}},
			wantErr: false,
		},
		{
			name:    "no follow",
			fields:  fields{ctx: context.Background()},
			args:    args{req: &relation.DouyinRelationFollowListRequest{UserId: 2, ToUserId: 2}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FollowListService{
				ctx: tt.fields.ctx,
			}
			got, err := s.FollowList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FollowList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
