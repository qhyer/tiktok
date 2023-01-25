package constants

const (
	JWTSigningKey = "8QpXHd59YYhk5IMj"

	UserTableName  = "user"
	VideoTableName = "video"

	ApiServiceName     = "api"
	UserServiceName    = "user"
	FeedServiceName    = "feed"
	PublishServiceName = "publish"

	MySQLDefaultDSN = "root:DFNFoTdxTfPY3B7X@tcp(103.200.115.51:3306)/tiktok?charset=utf8&parseTime=True&loc=Local"
	VideoQueryLimit = 30

	EtcdAddress = "127.0.0.1:2379"

	CPURateLimit float64 = 80.0

	OSSBaseURL         = "https://tt-test.qhyer.com/assets/"
	VideoBucketName    = "video"
	CoverBucketName    = "cover"
	OSSEndPoint        = "127.0.0.1:9000"
	OSSAccessKeyID     = ""
	OSSSecretAccessKey = ""
	OSSDefaultExpiry   = 3600
)
