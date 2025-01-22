package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wqh/easy/distribute/registry"
)

func main() {

	registry.SetUpRegisterService() // 启动心跳检测

	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server

	srv.Addr = registry.ServicePort

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Printf("%s shutdown successfully\n", "registry service")
			} else {
				log.Printf("%s shutdown failed, err: %v\n", "registry service", err)
			}
		}
		cancel()
	}()

	go func() {
		log.Println("registry service start up. press ctrl + c to stop it")
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("%s shutdown failed, err: %v\n", "registry service", err)
		}
		cancel()
	}()
	<-ctx.Done()
	log.Println("registry service terminated")
}
