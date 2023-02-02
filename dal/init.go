package dal

import (
	"tiktok/dal/mysql"
	"tiktok/dal/neo4j"
	"tiktok/dal/redis"
)

// Init init dal
func Init() {
	mysql.Init() // mysql init
	neo4j.Init() // neo4j init
	redis.Init() // redis init
}
