package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"topcoint/pkg/config"
	"topcoint/types"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *CurrencyInfo) GetCurrencyInfo(w http.ResponseWriter, r *http.Request) {
	currencyInfoResponse := types.CurrencyInfoResponse{Data: make(map[string]interface{})}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		conn.WriteJSON(map[string]string{"error": "Missing symbol parameter " + err.Error()})
		return
	}

	pagination := r.URL.Query().Get("pagination")
	if pagination == "" {
		conn.WriteJSON(map[string]string{"error": "Missing pagination parameter " + err.Error()})
		return
	}

	page := r.URL.Query().Get("page")

	req := types.CurrencyInfoRequest{
		Symbol:     symbol,
		Pagination: pagination,
	}

	currencyInfo := types.CurrencyInfoAPIResponse{}
	queryCurrencyInfo := fmt.Sprintf(c.cfg.ApiURL+"/asset/v2/metadata?asset_lookup_priority=SYMBOL&assets=%s&asset_language=en-US&quote_asset=USD", req.Symbol)
	err = fetchAPIData(c.cfg.ApiKey, queryCurrencyInfo, &currencyInfo)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Error fetching currency info data " + err.Error()})
		return
	}

	currencyStats := types.CurrencyStatsAPIResponse{}
	queryCurrencyStats := fmt.Sprintf(c.cfg.ApiURL+"/asset/v1/top/list?page=%s&page_size=%s&sort_by=CIRCULATING_MKT_CAP_USD&sort_direction=DESC&groups=ID,BASIC,SUPPLY,PRICE,MKT_CAP,VOLUME,CHANGE,TOPLIST_RANK&toplist_quote_asset=USD", page, req.Pagination)
	err = fetchAPIData(c.cfg.ApiKey, queryCurrencyStats, &currencyStats)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Error fetching currency stats data " + err.Error()})
		return
	}

	for k, v := range currencyInfo.Data {
		currencyInfoResponse.Data[k] = fmt.Sprintf("currencyInfo %v", v)
	}

	for k, v := range currencyStats.Data {
		currencyInfoResponse.Data[k] = fmt.Sprintf("currencyStats %v", v)
	}

	conn.WriteJSON(currencyInfoResponse)
}

func fetchAPIData(apiKey, url string, dest interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Set("authorization", "Apikey "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp types.APIResponseWithError
	if err = json.Unmarshal(body, &apiResp); err == nil && apiResp.Err != nil && apiResp.Err.Message != "" {
		fmt.Println("apiResp.Err.Message", apiResp.Err.Message)
		return fmt.Errorf("API error: %s", apiResp.Err.Message)
	}

	if err = json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
