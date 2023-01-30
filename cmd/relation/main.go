package main

import (
	"net"

	"tiktok/dal"
	relation "tiktok/kitex_gen/relation/relationsrv"
	"tiktok/pkg/bound"
	"tiktok/pkg/constants"
	"tiktok/pkg/middleware"
	"tiktok/pkg/minio"
	"tiktok/pkg/rpc"
	tracer2 "tiktok/pkg/tracer"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

func Init() {
	tracer2.InitJaeger(constants.RelationServiceName)
	dal.Init()
	rpc.InitFeedRpc()
	minio.Init()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8894")
	if err != nil {
		panic(err)
	}
	Init()
	svr := relation.NewServer(new(RelationSrvImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.RelationServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                                 // middleware
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