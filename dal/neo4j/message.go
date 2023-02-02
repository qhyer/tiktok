package neo4j

import (
	"context"

	"tiktok/kitex_gen/message"
	"tiktok/pkg/errno"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// UpsertLastMessage upsert last message
func UpsertLastMessage(ctx context.Context, userId int64, toUserId int64, content string) (success bool, err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close(ctx)
	}()
	res, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := isFriendRelation(ctx, tx, userId, toUserId)
		if err != nil {
			return false, err
		}
		res, err = upsertLastMessage(ctx, tx, userId, toUserId, content)
		if err != nil {
			return false, err
		}
		return res, nil
	})
	if err != nil {
		return false, err
	}
	success = res.(bool)
	return success, nil
}

// MQueryLastMessage query friends last message
func MQueryLastMessage(ctx context.Context, userId int64, toUserIds []int64) (messages []*message.Message, err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close(ctx)
	}()
	res, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := mQueryLastMessage(ctx, tx, userId, toUserIds)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	messages = res.([]*message.Message)
	return messages, nil
}

func mQueryLastMessage(ctx context.Context, tx neo4j.ManagedTransaction, userId int64, toUserIds []int64) (messages []*message.Message, err error) {
	result, err := tx.Run(ctx,
		"MATCH (a:User{id: $userId})-[m:LastMessage]-(b) WHERE b.id IN $toUserIds RETURN {id:a.id, to_id:b.id, content:m.content, sender:m.sender} AS b;",
		map[string]interface{}{
			"userId":    userId,
			"toUserIds": toUserIds,
		},
	)
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	messages = make([]*message.Message, 0, len(records))
	for _, r := range records {
		// 获取好友
		u, ok := r.Get("b")
		if !ok {
			continue
		}
		node := u.(map[string]any)
		uid := node["id"].(int64)
		if uid == 0 {
			continue
		}
		tuid := node["to_id"].(int64)
		if tuid == 0 {
			continue
		}
		msg := node["content"].(string)
		if msg == "" {
			continue
		}
		sender := node["sender"].(int64)
		if sender == uid {
			messages = append(messages, &message.Message{
				FromUserId: uid,
				ToUserId:   tuid,
				Content:    msg,
			})
		} else {
			messages = append(messages, &message.Message{
				FromUserId: tuid,
				ToUserId:   uid,
				Content:    msg,
			})
		}
	}
	return messages, nil
}

func upsertLastMessage(ctx context.Context, tx neo4j.ManagedTransaction, userId int64, toUserId int64, content string) (bool, error) {
	query := "MATCH (a:User{id: $userId}), (b:User{id: $toUserId}) MERGE (a)-[m:LastMessage]-(b) SET m = {content:$content, sender:$userId} RETURN m;"
	parameters := map[string]interface{}{
		"userId":   userId,
		"toUserId": toUserId,
		"content":  content,
	}
	res, err := tx.Run(ctx, query, parameters)
	if err != nil {
		return false, err
	}
	record, err := res.Single(ctx)
	if err != nil {
		return false, err
	}
	_, found := record.Get("m")
	if !found {
		return false, errno.DBOperationFailedErr
	}
	return true, err
}
