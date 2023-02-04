package constants

import "time"

// 配置信息
const (
	RedisAddress  = "119.91.157.116:6379"
	RedisPassword = "uEmf1RH07o6NiMVg"
)

// 缓存key

const (
	RedisUserKey        = "USER:%d"
	RedisVideoKey       = "VIDEO:%d"
	RedisFeedKey        = "FEED"
	RedisCommentListKey = "COMMENT_LIST:%d"
)

// 超时时间

const (
	MaxRandExpireSecond = 300
	UserInfoExpiry      = 10 * time.Minute
	VideoInfoExpiry     = 10 * time.Minute
)
