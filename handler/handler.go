package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"topcoint/pkg/config"
)

type CryptoCurrencier interface {
	GetCryptoCurrencies(w http.ResponseWriter, r *http.Request)
}

type CryptoCurrencies struct {
	cfg config.Config
}

func NewCryptoCurrencies(cfg config.Config) CryptoCurrencier {
	return &CryptoCurrencies{cfg: cfg}
}

func (c *CryptoCurrencies) GetCryptoCurrencies(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var cd json.RawMessage

	symbol := r.URL.Query().Get("symbol")
	url := "http://" + c.cfg.Hostname + ":" + c.cfg.ApiPort + "/" + c.cfg.ApiServiceName
	if symbol != "" {
		url += "?symbol=" + symbol
	}

	wg.Add(1)
	go fetch(url, &cd, &wg)
	wg.Wait()

	response := map[string]json.RawMessage{
		"coindesk": cd,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func fetch(url string, result *json.RawMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		*result = json.RawMessage(`null`)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		*result = json.RawMessage(`null`)
		return
	}

	*result = body
}
