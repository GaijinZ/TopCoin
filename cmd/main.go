package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"topcoint/pkg/config"
	"topcoint/pkg/handler"
	"topcoint/pkg/router"
	"topcoint/pkg/server"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	wg := sync.WaitGroup{}
	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	currencyInfo := handler.NewCurrencyInfo(*cfg)

	r := router.Router(currencyInfo)

	srv, err := server.NewServer(
		server.WithRouter(r),
		server.WithHost(cfg.HostName),
		server.WithPort(cfg.ApiPort),
	)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	shutdownSignalChan := make(chan struct{})

	waitForShutdownTrigger := func(chw chan struct{}, wg *sync.WaitGroup) {
		wg.Wait()
		chw <- struct{}{}
		close(chw)
	}

	wg.Add(1)
	go func() {
		srv.Run()
		wg.Done()
	}()

	go waitForShutdownTrigger(shutdownSignalChan, &wg)

	select {
	case <-shutdownSignalChan:
		fmt.Println("Received a shutdown signal, shutting down gracefully...")
	case <-interrupt:
		fmt.Println("Received a shutdown signal...")
	case <-ctx.Done():
		fmt.Println("Context cancelled, shutting down...")
	}
}
