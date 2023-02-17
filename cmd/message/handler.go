package main

import (
	"context"

	"tiktok/cmd/message/service"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/message"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
)

// MessageSrvImpl implements the last service interface defined in the IDL.
type MessageSrvImpl struct{}

// MessageAction implements the MessageSrvImpl interface.
func (s *MessageSrvImpl) MessageAction(ctx context.Context, req *message.DouyinMessageActionRequest) (resp *message.DouyinMessageActionResponse, err error) {
	resp = new(message.DouyinMessageActionResponse)

	switch req.GetActionType() {
	case constants.SendMessageAction:
		err = service.NewMessageActionService(ctx).SendMessage(req)
	default:
		err = errno.ParamErr
	}
	if err != nil {
		resp = pack.BuildMessageActionResp(err)
		return resp, nil
	}

	resp = pack.BuildMessageActionResp(errno.Success)
	return resp, nil
}

// MessageList implements the MessageSrvImpl interface.
func (s *MessageSrvImpl) MessageList(ctx context.Context, req *message.DouyinMessageListRequest) (resp *message.DouyinMessageListResponse, err error) {
	resp = new(message.DouyinMessageListResponse)

	messages, err := service.NewMessageListService(ctx).MessageList(req)
	if err != nil {
		resp = pack.BuildMessageListResp(err)
		return resp, nil
	}

	resp = pack.BuildMessageListResp(errno.Success)
	resp.MessageList = messages
	return resp, nil
}
