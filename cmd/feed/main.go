package main

import (
	"net"
	"tiktok/cmd/user/dal"
	"tiktok/pkg/bound"
	"tiktok/pkg/constants"
	"tiktok/pkg/middleware"

	feed "tiktok/kitex_gen/feed/feedsrv"
	tracer2 "tiktok/pkg/tracer"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

func Init() {
	tracer2.InitJaeger(constants.FeedServiceName)
	dal.Init()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8890")
	if err != nil {
		panic(err)
	}
	Init()
	svr := feed.NewServer(new(FeedSrvImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.UserServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                             // middleware
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr), // address
		//server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                           // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),     // tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()), // BoundHandler
		server.WithRegistry(r),                              // registry
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
