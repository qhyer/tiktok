package service

import (
	"context"
	"log"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/relation"
)

func TestQueryRelationService_IsFriend(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *relation.DouyinRelationIsFriendRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{ctx: context.Background()},
			args: args{&relation.DouyinRelationIsFriendRequest{
				UserId:   8,
				ToUserId: 10,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &QueryRelationService{
				ctx: tt.fields.ctx,
			}
			got, err := s.IsFriend(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFriend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
