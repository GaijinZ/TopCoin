package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"topcoint/coindesk/handler"
	"topcoint/coindesk/pkg/config"
	"topcoint/coindesk/router"
	"topcoint/coindesk/server"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	currencyInfo := handler.NewCurrencyInfo(*cfg)

	listener, err := net.Listen("tcp", ":"+cfg.ApiPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}

	r := router.Router(currencyInfo)

	srv := server.NewServer(listener, r)

	go srv.Run()

	select {
	case <-interrupt:
		fmt.Println("Received a shutdown signal...")
	case <-ctx.Done():
		fmt.Println("Context cancelled, shutting down...")
	}
}
