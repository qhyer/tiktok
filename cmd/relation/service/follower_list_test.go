package service

import (
	"context"
	"log"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
)

func TestFollowerListService_FollowerList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *relation.DouyinRelationFollowerListRequest
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
			name:    "no follower",
			fields:  fields{ctx: context.Background()},
			args:    args{req: &relation.DouyinRelationFollowerListRequest{UserId: 1}},
			wantErr: false,
		},
		{
			name:    "ok",
			fields:  fields{ctx: context.Background()},
			args:    args{req: &relation.DouyinRelationFollowerListRequest{UserId: 2}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FollowerListService{
				ctx: tt.fields.ctx,
			}
			got, err := s.FollowerList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FollowerList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
