package service

import (
	"context"
	"testing"

	"tiktok/dal/db"
	"tiktok/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

func TestCheckUser(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		testName string
		args     args
		wantErr  bool
	}{
		{
			testName: "ok",
			args: args{
				username: "qhyer",
				password: "qhyer",
			},
			wantErr: false,
		},
		{
			testName: "user not exists",
			args: args{
				username: "q2",
				password: "234",
			},
			wantErr: true,
		},
		{
			testName: "wrong password",
			args: args{
				username: "qhyer",
				password: "123456",
			},
			wantErr: true,
		},
	}
	db.Init()
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewCheckUserService(context.Background()).CheckUser(&user.DouyinUserLoginRequest{
				Username: tt.args.username,
				Password: tt.args.password,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			klog.Info(tt.testName + " success")
		})
	}
}
