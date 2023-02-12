package pack

import (
	"tiktok/dal/mysql"
	"tiktok/kitex_gen/message"
)

// Message pack message
func Message(msg *mysql.Message) *message.Message {
	if msg == nil {
		return nil
	}
	return &message.Message{
		Id:         msg.Id,
		ToUserId:   msg.ToUserId,
		FromUserId: msg.UserId,
		Content:    msg.Content,
		CreateTime: msg.CreatedAt.UnixMilli(),
	}
}

// Messages pack list of message
func Messages(msgs []*mysql.Message) []*message.Message {
	messages := make([]*message.Message, 0, len(msgs))
	if len(msgs) == 0 {
		return messages
	}

	// pack message
	for _, m := range msgs {
		if msg := Message(m); msg != nil {
			messages = append(messages, msg)
		}
	}

	return messages
}
