package neo4j

import (
	"context"

	"tiktok/kitex_gen/user"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateUser(ctx context.Context, user *user.User) (err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err := session.Close(ctx)
		if err != nil {
			return
		}
	}()
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := "CREATE (:User {id: $id, username: $username, follow_count: 0, follower_count: 0})"
		parameters := map[string]interface{}{
			"id":       user.Id,
			"username": user.Name,
		}
		_, err = tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func MGetUserByUserIds(ctx context.Context, userIds []int64) (users []*user.User, err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err := session.Close(ctx)
		if err != nil {
			return
		}
	}()
	res, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return queryUserInfoByUserIds(ctx, tx, userIds)
	})
	if err != nil {
		return nil, err
	}
	users = res.([]*user.User)
	return users, nil
}

func queryUserInfoByUserIds(ctx context.Context, tx neo4j.ManagedTransaction, userIds []int64) ([]*user.User, error) {
	result, err := tx.Run(ctx,
		"MATCH(a:User) WHERE a.id IN $userIds RETURN a;",
		map[string]interface{}{
			"userIds": userIds,
		},
	)
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*user.User, 0, len(records))
	for _, r := range records {
		u, ok := r.Get("a")
		if !ok {
			continue
		}
		node := u.(neo4j.Node)
		uid, err := neo4j.GetProperty[int64](node, "id")
		if err != nil {
			continue
		}
		username, err := neo4j.GetProperty[string](node, "username")
		if err != nil {
			continue
		}
		followCount, err := neo4j.GetProperty[int64](node, "follow_count")
		if err != nil {
			continue
		}
		followerCount, err := neo4j.GetProperty[int64](node, "follower_count")
		if err != nil {
			continue
		}
		users = append(users, &user.User{
			Id:            uid,
			Name:          username,
			FollowCount:   &followCount,
			FollowerCount: &followerCount,
			IsFollow:      false,
		})
	}
	return users, nil
}
