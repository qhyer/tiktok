package pack

import (
	"tiktok/dal/mysql"
	"tiktok/kitex_gen/user"
)

// User pack user info
func User(u *mysql.User) *user.User {
	if u == nil {
		return nil
	}

	return &user.User{
		Id:   u.Id,
		Name: u.UserName,
	}
}

// Users pack list of user info
func Users(us []*mysql.User) []*user.User {
	users := make([]*user.User, 0)
	for _, u := range us {
		if user2 := User(u); user2 != nil {
			users = append(users, user2)
		}
	}
	return users
}
