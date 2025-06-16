package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"topcoint/coindesk/pkg/config"
)

type Currencier interface {
	GetCurrencyInfo(w http.ResponseWriter, r *http.Request)
}

type CurrencyInfo struct {
	cfg config.Config
}

func NewCurrencyInfo(cfg config.Config) Currencier {
	return &CurrencyInfo{cfg: cfg}
}

func (c *CurrencyInfo) GetCurrencyInfo(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Missing symbol parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf(c.cfg.ApiURL+"/asset/v2/metadata?asset_lookup_priority=SYMBOL&assets=%s&asset_language=en-US&quote_asset=USD", symbol)
	currencyInfo, err := fetchAPIData(c.cfg.ApiKey, url)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	url2 := c.cfg.ApiURL + "/asset/v1/top/list?page=1&page_size=10&sort_by=CIRCULATING_MKT_CAP_USD&sort_direction=DESC&groups=ID,BASIC,SUPPLY,PRICE,MKT_CAP,VOLUME,CHANGE,TOPLIST_RANK&toplist_quote_asset=USD"
	currencyStats, err := fetchAPIData(c.cfg.ApiKey, url2)
	if err != nil {
		http.Error(w, "Error fetching data from API2", http.StatusInternalServerError)
		return
	}

	combined := map[string]interface{}{
		"currencyInfo":  currencyInfo,
		"currencyStats": currencyStats,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(combined); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func fetchAPIData(apiKey, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Set("authorization", "Apikey "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}
