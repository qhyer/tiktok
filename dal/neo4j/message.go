package neo4j

import (
	"context"
	"log"

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
		"MATCH (a:User{id: $userId})-[m:LastMessage]-(b) WHERE b.id IN $toUserIds RETURN {id:a.id, to_id:b.id, content:m.content, is_sender:m.is_sender};",
		map[string]interface{}{
			"userId":    userId,
			"toUserIds": toUserIds,
		},
	)
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	log.Print(records)
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
		node := u.(neo4j.Node)
		uid, err := neo4j.GetProperty[int64](node, "id")
		if err != nil {
			continue
		}
		msg, err := neo4j.GetProperty[string](node, "content")
		if err != nil {
			continue
		}
		isSender, err := neo4j.GetProperty[int64](node, "is_sender")
		if err != nil {
			continue
		}
		if isSender == 1 {
			messages = append(messages, &message.Message{
				FromUserId: userId,
				ToUserId:   uid,
				Content:    msg,
			})
		} else {
			messages = append(messages, &message.Message{
				FromUserId: uid,
				ToUserId:   userId,
				Content:    msg,
			})
		}
	}
	return messages, nil
}

func upsertLastMessage(ctx context.Context, tx neo4j.ManagedTransaction, uid1 int64, uid2 int64, content string) (bool, error) {
	query := "MATCH (a:User{id: $uid1}), (b:User{id: $uid2}) MERGE (a)-[m:LastMessage]-(b) MERGE (b)-[n:LastMessage]-(a) SET m = {content:$content, is_sender:1} SET n = {content:$content, is_sender:0} RETURN m, n;"
	parameters := map[string]interface{}{
		"uid1":    uid1,
		"uid2":    uid2,
		"content": content,
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
	_, found = record.Get("n")
	if !found {
		return false, errno.DBOperationFailedErr
	}
	return true, err
}
