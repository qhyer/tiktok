package main

import (
	"net"

	"tiktok/cmd/rpc"
	"tiktok/dal"
	message "tiktok/kitex_gen/message/messagesrv"
	"tiktok/pkg/bound"
	"tiktok/pkg/censor"
	"tiktok/pkg/constants"
	"tiktok/pkg/middleware"
	"tiktok/pkg/tracer"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	trace "github.com/kitex-contrib/tracer-opentracing"
	nacos "github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func Init() {
	tracer.InitJaeger(constants.MessageServiceName)
	dal.Init()
	rpc.InitUserRpc()
	censor.Init()
}

func main() {
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
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8894")
	if err != nil {
		panic(err)
	}
	Init()
	svr := message.NewServer(new(MessageSrvImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.MessageServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                                // middleware
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr),                                       // address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),                    // tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()),                // BoundHandler
		server.WithRegistry(registry.NewNacosRegistry(cli)),                // registry
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
