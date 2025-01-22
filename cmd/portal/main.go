package main

import (
	"context"
	"fmt"
	"github.com/wqh/easy/distribute/portal"
	"github.com/wqh/easy/distribute/registry"
	"github.com/wqh/easy/distribute/service"
	stlog "log"
)

/**
* description:
* author: wqh
* date: 2025/1/21
 */

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}

	host, port := "localhost", "5681"
	address := fmt.Sprintf("http://%s:%s", host, port)
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		registry.Registration{
			ServiceName: registry.PortalService,
			ServiceURL:  address,
			RequiredServices: []registry.ServiceName{
				registry.GradeService,
				registry.LogService,
			},
			ServiceUpdateURL: fmt.Sprintf("%s/services", address),
			HeartbeatURL:     fmt.Sprintf("%s/heartbeat", address),
		},
		portal.RegisterHandlerFunc,
	)
	if err != nil {
		stlog.Fatalf("service start failed, err: %v\n", err)
	}
	<-ctx.Done()
}
