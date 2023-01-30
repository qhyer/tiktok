package main

import (
	"context"

	"tiktok/kitex_gen/message"
)

// MessageSrvImpl implements the last service interface defined in the IDL.
type MessageSrvImpl struct{}

// MessageAction implements the MessageSrvImpl interface.
func (s *MessageSrvImpl) MessageAction(ctx context.Context, req *message.DouyinMessageActionRequest) (resp *message.DouyinMessageActionResponse, err error) {
	// TODO: Your code here...
	return
}

// ChatList implements the MessageSrvImpl interface.
func (s *MessageSrvImpl) ChatList(ctx context.Context, req *message.DouyinMessageChatRequest) (resp *message.DouyinMessageChatResponse, err error) {
	// TODO: Your code here...
	return
}
