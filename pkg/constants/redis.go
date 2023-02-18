package constants

import "time"

// 配置信息
const (
	RedisAddress  = "119.91.157.116:6379"
	RedisPassword = "uEmf1RH07o6NiMVg"
)

// 缓存key

const (
	RedisUserKey         = "USER:%d"
	RedisVideoKey        = "VIDEO:%d"
	RedisFeedKey         = "FEED"
	RedisCommentKey      = "COMMENT:%d"
	RedisCommentListKey  = "COMMENT_LIST:%d"
	RedisPublishListKey  = "PUBLISH_LIST:%d"
	RedisFavoriteListKey = "FAVORITE_LIST:%d"
	RedisFollowListKey   = "FOLLOW_LIST:%d"
	RedisFollowerListKey = "FOLLOWER_LIST:%d"
	RedisFriendListKey   = "FRIEND_LIST:%d"
	RedisMessageListKey  = "MESSAGE_LIST:%d:%d"
	RedisMessageKey      = "MESSAGE:%d"
)

// 超时时间

const (
	MaxRandExpireSecond = 300
	UserInfoExpiry      = 10 * time.Minute
	VideoInfoExpiry     = 10 * time.Minute
	CommentExpiry       = 15 * time.Minute
	CommentListExpiry   = 5 * time.Minute
	PublishListExpiry   = 15 * time.Minute
	FavoriteListExpiry  = 5 * time.Minute
	FollowListExpiry    = 5 * time.Minute
	FollowerListExpiry  = 5 * time.Minute
	FriendListExpiry    = 5 * time.Minute
	MessageListExpiry   = 10 * time.Minute
	MessageExpiry       = 10 * time.Minute
)
