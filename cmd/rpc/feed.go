package rpc

import (
	"context"
	"time"

	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/feed/feedsrv"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/middleware"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/kitex-contrib/registry-nacos/resolver"
	trace "github.com/kitex-contrib/tracer-opentracing"
	nacos "github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var feedClient feedsrv.Client

func InitFeedRpc() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(constants.NacosAddress, constants.NacosPort),
	}
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		Username:            constants.NacosUsername,
		Password:            constants.NacosPassword,
	}
	cli, err := nacos.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	c, err := feedsrv.NewClient(
		constants.FeedServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(100),                       // mux
		client.WithRPCTimeout(10*time.Second),               // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),      // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),   // retry
		client.WithSuite(trace.NewDefaultClientSuite()),     // tracer
		client.WithResolver(resolver.NewNacosResolver(cli)), // resolver
	)
	if err != nil {
		panic(err)
	}
	feedClient = c
}

func Feed(ctx context.Context, req *feed.DouyinFeedRequest) (*feed.DouyinFeedResponse, error) {
	resp, err := feedClient.Feed(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func GetVideosByVideoIdsAndCurrentUserId(ctx context.Context, req *feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest) (*feed.DouyinGetVideosByVideoIdsAndCurrentUserIdResponse, error) {
	resp, err := feedClient.GetVideosByVideoIdsAndCurrentUserId(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}
