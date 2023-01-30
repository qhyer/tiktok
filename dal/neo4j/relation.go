package neo4j

import (
	"context"
	"log"

	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// FollowAction user follow action
func FollowAction(ctx context.Context, userId int64, toUserId int64) (err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close(ctx)
	}()
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err = updateFollowNum(ctx, tx, userId, 1)
		if err != nil {
			return nil, err
		}
		_, err = updateFollowerNum(ctx, tx, toUserId, 1)
		if err != nil {
			return nil, err
		}
		_, err = addFollow(ctx, tx, userId, toUserId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

// UnfollowAction user unfollow action
func UnfollowAction(ctx context.Context, userId int64, toUserId int64) (err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close(ctx)
	}()
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err = updateFollowNum(ctx, tx, userId, -1)
		if err != nil {
			return nil, err
		}
		_, err = updateFollowerNum(ctx, tx, toUserId, -1)
		if err != nil {
			return nil, err
		}
		_, err = deleteFollow(ctx, tx, userId, toUserId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

// FollowList get user follow list
func FollowList(ctx context.Context, userId int64) (users []*user.User, err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close(ctx)
	}()
	res, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		users, err := queryUserFollow(ctx, tx, userId)
		if err != nil {
			return nil, err
		}
		return users, nil
	})
	if err != nil {
		return nil, err
	}
	users = res.([]*user.User)
	return users, nil
}

// FollowerList get user follower list
func FollowerList(ctx context.Context, userId int64) (users []*user.User, err error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close(ctx)
	}()
	res, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		users, err := queryUserFollower(ctx, tx, userId)
		if err != nil {
			return nil, err
		}
		return users, nil
	})
	if err != nil {
		return nil, err
	}
	users = res.([]*user.User)
	return users, nil
}

func updateFollowNum(ctx context.Context, tx neo4j.ManagedTransaction, userId int64, addNum int64) (interface{}, error) {
	query := "MATCH (a:User {id: $userId}) " +
		"SET a.follow_count = a.follow_count + $addNum RETURN a;"
	parameters := map[string]interface{}{
		"userId": userId,
		"addNum": addNum,
	}
	res, err := tx.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}
	record, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}
	_, found := record.Get("a")
	if !found {
		return nil, errno.DBOperationFailedErr
	}
	return res, err
}

func updateFollowerNum(ctx context.Context, tx neo4j.ManagedTransaction, userId int64, addNum int64) (interface{}, error) {
	query := "MATCH (a:User {id: $userId}) " +
		"SET a.follower_count = a.follower_count + $addNum RETURN a;"
	parameters := map[string]interface{}{
		"userId": userId,
		"addNum": addNum,
	}
	res, err := tx.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}
	record, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}
	_, found := record.Get("a")
	if !found {
		return nil, errno.DBOperationFailedErr
	}
	return res, err
}

func addFollow(ctx context.Context, tx neo4j.ManagedTransaction, userId int64, toUserId int64) (interface{}, error) {
	query := "MATCH (a:User), (b:User) " +
		"WHERE a.id = $userId AND b.id = $toUserId AND NOT (a)-[:Follow]->(b) " +
		"CREATE (a)-[r:Follow]->(b) RETURN r;"
	parameters := map[string]interface{}{
		"userId":   userId,
		"toUserId": toUserId,
	}
	res, err := tx.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}
	record, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}
	ret, found := record.Get("r")
	log.Println(ret)
	if !found {
		return nil, errno.DBOperationFailedErr
	}
	return res, err
}

func deleteFollow(ctx context.Context, tx neo4j.ManagedTransaction, userId int64, toUserId int64) (interface{}, error) {
	query := "MATCH (a:User)-[r:Follow]->(b:User) " +
		"WHERE a.id = $userId AND b.id = $toUserId " +
		"DELETE r RETURN r;"
	parameters := map[string]interface{}{
		"userId":   userId,
		"toUserId": toUserId,
	}
	res, err := tx.Run(ctx, query, parameters)
	if err != nil {
		return nil, err
	}
	record, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}
	_, found := record.Get("r")
	if !found {
		return nil, errno.DBOperationFailedErr
	}
	return res, err
}

func queryUserFollow(ctx context.Context, tx neo4j.ManagedTransaction, userId int64) ([]*user.User, error) {
	result, err := tx.Run(ctx,
		"MATCH (:User {id: $userId})-[r:Follow]->(User) "+
			"RETURN User;",
		map[string]interface{}{
			"userId": userId,
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
		u, ok := r.Get("User")
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
			IsFollow:      true,
		})
	}
	return users, nil
}

func queryUserFollower(ctx context.Context, tx neo4j.ManagedTransaction, userId int64) ([]*user.User, error) {
	result, err := tx.Run(ctx,
		"MATCH (User)-[r:Follow]->(:User{id: $userId}) "+
			"RETURN User;",
		map[string]interface{}{
			"userId": userId,
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
		u, ok := r.Get("User")
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
		})
	}
	return users, nil
}
