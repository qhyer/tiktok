package initialize

import (
	"tiktok/cmd/rpc"
)

func Rpc() {
	rpc.InitUserRpc()
	rpc.InitFeedRpc()
	rpc.InitPublishRpc()
	rpc.InitFavoriteRpc()
	rpc.InitCommentRpc()
	rpc.InitRelationRpc()
	rpc.InitMessageRpc()
}
