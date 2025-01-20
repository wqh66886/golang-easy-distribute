package service

import (
	"context"
	"fmt"
	"github.com/wqh/easy/distribute/registry"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Start The entry for service startup
func Start(ctx context.Context, host, port string, reg registry.Registration, registerFunc func()) (context.Context, error) {
	// registration service take precede
	registerFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	// when service start up register service
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

// startService this function is used to start the service
func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				err = registry.DeregisterService(registry.Registration{
					ServiceName: serviceName,
					ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
				})
				if err != nil {
					log.Printf("%s remove failed, err : %v\n", serviceName, err)
				}
			} else {
				log.Printf("%s service is shutdown", serviceName)
			}
		}
		cancel()
	}()

	go func() {
		log.Printf("%s is started up on port %s. Press ctl + c to stop it", serviceName, port)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		// listen signal to remove service from registry
		err := registry.DeregisterService(registry.Registration{
			ServiceName: serviceName,
			ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
		})
		if err != nil {
			log.Printf("%s remove failed, err : %v\n", serviceName, err)
		}
		err = srv.Shutdown(ctx)
		if err != nil {
			log.Printf("%s shutdown failed, err: %v\n", serviceName, err)
		}
		cancel()
	}()

	return ctx
}
