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
		err = session.Close(ctx)
	}()
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := "CREATE (:User {id: $id, username: $username, follow_count: 0, follower_count: 0})"
		parameters := map[string]interface{}{
			"id":       user.Id,
			"username": user.Name,
		}
		_, err = tx.Run(ctx, query, parameters)
		return nil, err
	}); err != nil {
		return err
	}
	return nil
}
