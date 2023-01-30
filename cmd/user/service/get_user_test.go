package service

import (
	"context"
	"log"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/user"
)

func TestMGetUserService_MGetUser(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *user.DouyinUserInfoRequest
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
			name: "qhyer",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				&user.DouyinUserInfoRequest{
					UserId:    1,
					ToUserIds: []int64{1, 2},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MGetUserService{
				ctx: tt.fields.ctx,
			}
			got, err := s.MGetUser(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
