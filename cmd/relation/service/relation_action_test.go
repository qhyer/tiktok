package service

import (
	"context"
	"testing"

	"tiktok/dal"
	"tiktok/kitex_gen/relation"
)

func TestFollowActionService_FollowAction(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *relation.DouyinRelationActionRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{ctx: context.Background()},
			args: args{&relation.DouyinRelationActionRequest{
				UserId:     1,
				ToUserId:   2,
				ActionType: 1,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationActionService{
				ctx: tt.fields.ctx,
			}
			if err := s.Follow(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFollowActionService_UnFollowAction(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		req *relation.DouyinRelationActionRequest
	}
	dal.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{ctx: context.Background()},
			args: args{&relation.DouyinRelationActionRequest{
				UserId:     1,
				ToUserId:   2,
				ActionType: 2,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationActionService{
				ctx: tt.fields.ctx,
			}
			if err := s.Unfollow(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Unfollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
