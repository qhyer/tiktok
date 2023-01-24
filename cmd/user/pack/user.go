package pack

import (
	"tiktok/cmd/user/dal/db"
	"tiktok/kitex_gen/user"
)

// User pack user info
func User(u *db.User) *user.User {
	if u == nil {
		return nil
	}

	followCount := u.FollowCount
	followerCount := u.FollowerCount

	return &user.User{
		Id:            u.Id,
		Name:          u.UserName,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      false,
	}
}

// Users pack list of user info
func Users(us []*db.User) []*user.User {
	users := make([]*user.User, 0)
	for _, u := range us {
		if user2 := User(u); user2 != nil {
			users = append(users, user2)
		}
	}
	return users
}
