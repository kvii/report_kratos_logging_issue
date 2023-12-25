package main

import (
	"context"
	"os"

	"github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"caller", log.DefaultCaller,
	)

	helper := log.NewHelper(logger)
	helper.Info("foo") // 保证 helper.Info 的 caller 没有问题

	// 等待客户端请求
	// 期望打印的 caller 是 logging.Server 打印 log 的位置
	srv := http.NewServer(
		http.Address(":9090"),
		http.Middleware(logging.Server(logger)),
	)
	helloworld.RegisterGreeterHTTPServer(srv, Service{})

	app := kratos.New(
		kratos.Logger(logger),
		kratos.Server(srv),
	)
	if err := app.Run(); err != nil {
		panic(err)
	}

	// INFO caller=server/server.go:20 msg=foo
	// INFO caller=http/server.go:317 msg=[HTTP] server listening on: [::]:9090
	// INFO caller=helloworld/helloworld_http.pb.go:62 kind=server component=http operation=/helloworld.Greeter/SayHello args=name:"k" code=0 reason= stack= latency=2.375e-05
}

type Service struct {
	helloworld.UnimplementedGreeterServer
}

func (Service) SayHello(context.Context, *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{}, nil
}
