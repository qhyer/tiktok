package service

import (
	"context"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/message"
)

func TestMessageActionService_SendMessage(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *message.DouyinMessageActionRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "not friend",
			fields: fields{context.Background()},
			args: args{&message.DouyinMessageActionRequest{
				UserId:     8,
				ToUserId:   9,
				ActionType: 1,
				Content:    "test",
			}},
			wantErr: true,
		},
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{&message.DouyinMessageActionRequest{
				UserId:     8,
				ToUserId:   10,
				ActionType: 1,
				Content:    "test",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MessageActionService{
				ctx: tt.fields.ctx,
			}
			if err := s.SendMessage(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
