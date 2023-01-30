package dal

import (
	"tiktok/dal/mysql"
	"tiktok/dal/neo4j"
)

// Init init dal
func Init() {
	mysql.Init() // mysql init
	neo4j.Init() // neo4j init
}
