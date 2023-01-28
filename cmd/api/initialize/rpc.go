package initialize

import "tiktok/pkg/rpc"

func Rpc() {
	rpc.InitUserRpc()
	rpc.InitFeedRpc()
	rpc.InitPublishRpc()
	rpc.InitFavoriteRpc()
}
