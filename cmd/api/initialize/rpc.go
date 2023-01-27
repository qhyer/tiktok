package initialize

import "tiktok/cmd/api/rpc"

func Rpc() {
	rpc.InitUserRpc()
	rpc.InitFeedRpc()
	rpc.InitPublishRpc()
	rpc.InitFavoriteRpc()
}
