package service

import (
	"context"
	"reflect"
	"testing"

	"tiktok/dal/db"
	"tiktok/kitex_gen/user"
)

func TestMGetUserService_MGetUser(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *user.DouyinUserInfoRequest
	}
	var zero int64
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
				ctx: nil,
			},
			args: args{
				&user.DouyinUserInfoRequest{
					UserId:    1,
					ToUserIds: []int64{1},
				},
			},
			want: []*user.User{
				{
					Id:            1,
					Name:          "qhyer",
					FollowerCount: &zero,
					FollowCount:   &zero,
				},
			},
			wantErr: false,
		},
	}
	db.Init()
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MGetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
