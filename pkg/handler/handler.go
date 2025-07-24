package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"topcoint/pkg/service"
	"topcoint/pkg/types"

	"github.com/gorilla/websocket"
	"topcoint/pkg/config"
)

type Currencier interface {
	Home(w http.ResponseWriter, r *http.Request)
	GetCurrencyInfo(w http.ResponseWriter, r *http.Request)
	CacheCurrency() error
	HandleSearch(w http.ResponseWriter, r *http.Request)
}

type CurrencyInfo struct {
	cfg   config.Config
	cache map[string]types.SummaryCryptoList
}

func NewCurrencyInfo(cfg config.Config) Currencier {
	return &CurrencyInfo{
		cfg:   cfg,
		cache: make(map[string]types.SummaryCryptoList),
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

	go service.HandleOutgoingMessages(conn, serverMessages)
	service.HandleIncomingMessages(c.cfg, conn, serverMessages)
}

func (c *CurrencyInfo) CacheCurrency() error {
	err := service.CacheCurrency(c.cfg, c.cache)
	if err != nil {
		return fmt.Errorf("failed to cache currency data: %w", err)
	}

	return nil
}

func (c *CurrencyInfo) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query param", http.StatusBadRequest)
		return
	}

	results := searchAssets(query, c.cache)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func searchAssets(symbol string, cache map[string]types.SummaryCryptoList) []types.CryptoCurrencyList {
	results := make([]types.CryptoCurrencyList, 0)
	query := strings.ToLower(symbol)

	for _, summary := range cache {
		for _, asset := range summary.Data.List {
			if strings.Contains(strings.ToLower(asset.Name), query) {
				results = append(results, asset)
			}
		}
	}

	return results
}
