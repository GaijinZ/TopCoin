package handler

import (
	"fmt"
	"net/http"

	"topcoint/pkg/config"
	"topcoint/pkg/service"
	"topcoint/pkg/types"

	"github.com/gorilla/websocket"
)

type Currencier interface {
	Home(w http.ResponseWriter, r *http.Request)
	APIClient(requests <-chan types.ClientMessage)
	WSHandler(requests chan<- types.ClientMessage) http.HandlerFunc
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

func (c *CurrencyInfo) APIClient(requests <-chan types.ClientMessage) {
	for req := range requests {
		stats, err := c.cache.HandleIncomingMessages(c.cfg, req)
		if err != nil {
			fmt.Println("Error fetching currency info:", err)
			continue
		}

		req.Reply <- types.CurrencyInfoResponse{Data: stats.Data}
	}
}

func (c *CurrencyInfo) WSHandler(requests chan<- types.ClientMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection:", err)
			return
		}
		defer conn.Close()

		for {
			var msg types.ClientMessage

			err = conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
				break
			}

			reply := make(chan types.CurrencyInfoResponse)
			msg.Reply = reply
			requests <- msg

			conn.WriteJSON(<-reply)
		}
	}
}

func (c *CurrencyInfo) CacheCurrency() error {
	err := c.cache.CacheCurrency(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to cache currency data: %w", err)
	}

	return nil
}
