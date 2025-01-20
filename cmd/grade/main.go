package main

import (
	"context"
	"fmt"
	"github.com/wqh/easy/distribute/grades"
	"github.com/wqh/easy/distribute/log"
	"github.com/wqh/easy/distribute/registry"
	"github.com/wqh/easy/distribute/service"
	stlog "log"
)

func main() {
	host, port := "localhost", "5680"
	address := fmt.Sprintf("http://%s:%s", host, port)
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		registry.Registration{
			ServiceName: registry.GradeService,
			ServiceURL:  address,
			RequiredServices: []registry.ServiceName{
				registry.LogService,
			},
			ServiceUpdateURL: fmt.Sprintf("%s/services", address),
		},
		grades.RegisterHandleFunc,
	)
	if err != nil {
		stlog.Fatalf("service start failed, err: %v\n", err)
	}

	if logProvider, err := registry.GetProviders(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider[0])
		log.SetClientLogger(logProvider[0], registry.GradeService)
	}

	<-ctx.Done()
	stlog.Println("shutting down grade service")
}
