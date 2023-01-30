package constants

import "time"

const (
	JWTSigningKey = "8QpXHd59YYhk5IMj"

	UserTableName     = "user"
	VideoTableName    = "video"
	FavoriteTableName = "favorite"
	CommentTableName  = "comment"

	ApiServiceName      = "api"
	UserServiceName     = "user"
	FeedServiceName     = "feed"
	PublishServiceName  = "publish"
	FavoriteServiceName = "favorite"
	CommentServiceName  = "comment"
	RelationServiceName = "relation"

	DoCommentAction     = 1
	DeleteCommentAction = 2

	DoFavoriteAction     = 1
	CancelFavoriteAction = 2

	DoFollowAction = 1
	UnfollowAction = 2

	MySQLDefaultDSN = "root:DFNFoTdxTfPY3B7X@tcp(119.91.157.116:3306)/tiktok?charset=utf8&parseTime=True&loc=Local"
	VideoQueryLimit = 30

	Neo4jDefaultURI = "neo4j://119.91.157.116:7687"
	Neo4jUser       = "neo4j"
	Neo4jPassword   = "MXOzml26024SyZl"

	EtcdAddress = "119.91.157.116:2379"

	CPURateLimit float64 = 80.0

	VideoBucketName    = "video"
	CoverBucketName    = "cover"
	OSSEndPoint        = "127.0.0.1:9000"
	OSSBaseUrl         = "https://tt-test.qhyer.com/assets"
	OSSAccessKeyID     = "tiktok"
	OSSSecretAccessKey = "ZqRNq8vd;9KLjx=9"
	OSSDefaultExpiry   = time.Hour
)
