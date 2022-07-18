package main

import (
	"flag"
	"fmt"

	"github.com/kevwan/grpc-restful/internal/config"
	"github.com/kevwan/grpc-restful/internal/server"
	"github.com/kevwan/grpc-restful/internal/svc"
	"github.com/kevwan/grpc-restful/pb"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sum.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	svr := server.NewSumServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSumServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	gw := gateway.MustNewServer(c.Gateway)
	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)
	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}
