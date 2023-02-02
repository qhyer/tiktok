package service

import (
	"context"
	"log"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/relation"
)

func TestFriendListService_FriendList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *relation.DouyinRelationFriendListRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*relation.FriendUser
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{context.Background()},
			args:    args{&relation.DouyinRelationFriendListRequest{UserId: 11}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FriendListService{
				ctx: tt.fields.ctx,
			}
			got, err := s.FriendList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FriendList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
