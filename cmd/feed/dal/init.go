package dal

import "tiktok/cmd/feed/dal/db"

// Init init dal
func Init() {
	db.Init() // mysql init
}
