package service

import (
	"context"
	"testing"
	"tiktok/cmd/user/dal/db"
	"tiktok/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

func TestRegister(t *testing.T) {
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
				username: "qhy",
				password: "qhy",
			},
			wantErr: false,
		},
		{
			testName: "username exists",
			args: args{
				username: "qhyer",
				password: "qhyer",
			},
			wantErr: true,
		},
	}
	db.Init()
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewCreateUserService(context.Background()).CreateUser(&user.DouyinUserRegisterRequest{
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
