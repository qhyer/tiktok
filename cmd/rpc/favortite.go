package rpc

import (
	"context"
	"time"

	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/favorite/favoritesrv"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/middleware"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

var favoriteClient favoritesrv.Client

func InitFavoriteRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := favoritesrv.NewClient(
		constants.FavoriteServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(100),                     // mux
		client.WithRPCTimeout(10*time.Second),             // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	favoriteClient = c
}

func FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (*favorite.DouyinFavoriteActionResponse, error) {
	resp, err := favoriteClient.FavoriteAction(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (*favorite.DouyinFavoriteListResponse, error) {
	resp, err := favoriteClient.FavoriteList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func GetUserFavoriteVideoIds(ctx context.Context, req *favorite.DouyinGetUserFavoriteVideoIdsRequest) (*favorite.DouyinGetUserFavoriteVideoIdsResponse, error) {
	resp, err := favoriteClient.GetUserFavoriteVideoIds(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, err
}