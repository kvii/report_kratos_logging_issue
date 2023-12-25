package main

import (
	"context"
	"os"

	"github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"caller", log.DefaultCaller,
	)

	helper := log.NewHelper(logger)
	helper.Info("bar") // 保证 helper.Info 的 caller 没有问题

	// 期望打印的 caller 是 logging.Client 打印 log 的位置
	ctx := context.Background()
	cli, err := http.NewClient(ctx,
		http.WithEndpoint("http://127.0.0.1:9090"),
		http.WithMiddleware(logging.Client(logger)),
	)
	if err != nil {
		panic(err)
	}

	// 发送 http 请求
	g := helloworld.NewGreeterHTTPClient(cli)
	_, err = g.SayHello(ctx, &helloworld.HelloRequest{Name: "k"})
	if err != nil {
		panic(err)
	}

	// INFO caller=client/client.go:19 msg=bar
	// INFO caller=http/client.go:277 kind=client component=http operation=/helloworld.Greeter/SayHello args=<nil> code=0 reason= stack= latency=0.003042
}
