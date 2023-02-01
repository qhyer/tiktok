package rpc

import "C"
import (
	"context"
	"time"

	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/relation/relationsrv"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/middleware"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

var relationClient relationsrv.Client

func InitRelationRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := relationsrv.NewClient(
		constants.RelationServiceName,
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
	relationClient = c
}

func RelationAction(ctx context.Context, req *relation.DouyinRelationActionRequest) (*relation.DouyinRelationActionResponse, error) {
	resp, err := relationClient.RelationAction(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func FollowList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (*relation.DouyinRelationFollowListResponse, error) {
	resp, err := relationClient.RelationFollowList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func FollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (*relation.DouyinRelationFollowerListResponse, error) {
	resp, err := relationClient.RelationFollowerList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func FriendList(ctx context.Context, req *relation.DouyinRelationFriendListRequest) (*relation.DouyinRelationFriendListResponse, error) {
	resp, err := relationClient.RelationFriendList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, err
}

func IsFriendRelation(ctx context.Context, req *relation.DouyinRelationIsFriendRequest) (*relation.DouyinRelationIsFriendResponse, error) {
	resp, err := relationClient.RelationIsFriend(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
