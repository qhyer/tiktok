package dal

import "tiktok/cmd/publish/dal/db"

// Init init dal
func Init() {
	db.Init() // mysql init
}
