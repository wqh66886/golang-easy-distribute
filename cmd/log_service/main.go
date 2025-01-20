package main

import (
	"context"
	"fmt"
	"github.com/wqh/easy/distribute/registry"
	stlog "log"

	"github.com/wqh/easy/distribute/log"
	"github.com/wqh/easy/distribute/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "5679"

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		registry.Registration{
			ServiceName: registry.LogService,
			ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
		},
		log.RegisterHandleFunc,
	)

	if err != nil {
		stlog.Fatalf("service start failed, err: %v\n", err)
	}
	<-ctx.Done()
	stlog.Println("shutting down log service")
}
