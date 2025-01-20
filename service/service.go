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

func Start(ctx context.Context, host, port string, reg registry.Registration, registerFunc func()) (context.Context, error) {
	registerFunc()
	// 启动log service 服务
	ctx = startService(ctx, reg.ServiceName, host, port)
	// 将log service 服务注册到注册中心
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("service start up failed, err: %v\n", err)
		}
		log.Println("service is stopped")
		// 从注册中心注销服务
		err = registry.DeregisterService(registry.Registration{
			ServiceName: serviceName,
			ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
		})
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		log.Printf("%s service is started up on %s. Press ctl + c to stop", serviceName, port)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		// 注销服务
		err := registry.DeregisterService(registry.Registration{
			ServiceName: serviceName,
			ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
		})
		if err != nil {
			log.Printf("remove service failed, err: %v\n", err)
		}
		err = srv.Shutdown(ctx)
		if err != nil {
			log.Printf("service shutdown failed, err: %v\n", err)
		} // 关闭服务
		cancel()
	}()

	return ctx
}
