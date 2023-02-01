package service

import (
	"context"
	"log"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/message"
)

func TestMessageListService_MessageList(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *message.DouyinMessageListRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*message.Message
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{context.Background()},
			args: args{&message.DouyinMessageListRequest{
				UserId:   9,
				ToUserId: 8,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MessageListService{
				ctx: tt.fields.ctx,
			}
			got, err := s.MessageList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}
