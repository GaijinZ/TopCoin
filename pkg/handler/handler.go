package handler

import (
	"fmt"
	"net/http"

	"topcoint/pkg/config"
	"topcoint/pkg/service"

	"github.com/gorilla/websocket"
)

type Currencier interface {
	Home(w http.ResponseWriter, r *http.Request)
	GetCurrencyInfo(w http.ResponseWriter, r *http.Request)
	CacheCurrency() error
}

type CurrencyInfo struct {
	cfg   config.Config
	cache service.Cacher
}

func NewCurrencyInfo(cfg config.Config, cache service.Cacher) Currencier {
	return &CurrencyInfo{
		cfg:   cfg,
		cache: cache,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *CurrencyInfo) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/home.html")
}

func (c *CurrencyInfo) GetCurrencyInfo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	serverMessages := make(chan interface{})
	defer close(serverMessages)

	go c.cache.HandleOutgoingMessages(conn, serverMessages)
	c.cache.HandleIncomingMessages(c.cfg, conn, serverMessages)
}

func (c *CurrencyInfo) CacheCurrency() error {
	err := c.cache.CacheCurrency(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to cache currency data: %w", err)
	}

	return nil
}
