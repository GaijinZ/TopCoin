package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"topcoint/pkg/config"
	"topcoint/pkg/types"

	"github.com/gorilla/websocket"
)

func HandleOutgoingMessages(conn *websocket.Conn, serverMessages <-chan interface{}) {
	for msg := range serverMessages {
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func HandleIncomingMessages(cfg config.Config, conn *websocket.Conn, serverMessages chan<- interface{}) {
	for {
		var msg types.ClientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			fmt.Println("Error reading JSON:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket closed unexpectedly: %v\n", err)
			}
			break
		}

		serverMessages <- map[string]string{"info": "Server-initiated message sent!"}

		if err := msg.Validate(); err != nil {
			conn.WriteJSON(map[string]string{"error": "Invalid input: " + err.Error()})
			return
		}

		response, err := buildCurrencyInfoResponse(cfg, msg)
		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			return
		}

		conn.WriteJSON(response)
	}
}

func CacheCurrency(cfg config.Config, cache map[string]types.SummaryCryptoList) error {
	url := fmt.Sprintf("%s/asset/v1/summary/list?asset_lookup_priority=SYMBOL", cfg.ApiURL)
	body, err := fetchAPIData(cfg.ApiKey, url)
	if err != nil {
		return fmt.Errorf("failed to fetch currency data: %w", err)
	}

	err = unmarshallCacheData(body, cache)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

func buildCurrencyInfoResponse(cfg config.Config, msg types.ClientMessage) (types.CurrencyInfoResponse, error) {
	response := types.CurrencyInfoResponse{Data: make(map[string]interface{})}

	if msg.Page == 1 {
		info, err := fetchCurrencyInfo(cfg, msg.Symbol)
		if err != nil {
			return response, fmt.Errorf("error fetching currency info: %w", err)
		}
		for k, v := range info {
			response.Data[k] = v
		}
	}

	stats, err := fetchCurrencyStats(cfg, msg.Page, msg.Pagination)
	if err != nil {
		return response, fmt.Errorf("error fetching currency stats: %w", err)
	}
	for i, v := range stats {
		response.Data[strconv.Itoa(i)] = v
	}

	return response, nil
}

func fetchCurrencyInfo(cfg config.Config, symbol string) (map[string]types.AssetInfo, error) {
	var res types.CurrencyInfoAPIResponse
	url := fmt.Sprintf("%s/asset/v2/metadata?asset_lookup_priority=SYMBOL&assets=%s&asset_language=en-US&quote_asset=USD", cfg.ApiURL, symbol)
	body, err := fetchAPIData(cfg.ApiKey, url)
	if err != nil {
		return nil, err
	}

	err = unmarshallData(body, &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %w", err)
	}

	for k, v := range res.Data {
		v.CreatedOnFormatted = time.Unix(v.CreatedOn, 0).Format(time.RFC3339)
		v.LunchDateFormatted = time.Unix(v.LunchDate, 0).Format(time.RFC3339)
		v.PriceLastUpdatedFormatted = time.Unix(v.PriceLastUpdated, 0).Format(time.RFC3339)
		res.Data[k] = v
	}

	return res.Data, nil
}

func fetchCurrencyStats(cfg config.Config, page, pagination int) (map[int]types.AssetInfo, error) {
	var res types.CurrencyStatsAPIResponse
	url := fmt.Sprintf("%s/asset/v1/top/list?page=%d&page_size=%d&sort_by=CIRCULATING_MKT_CAP_USD&sort_direction=DESC&groups=ID,BASIC,SUPPLY,PRICE,MKT_CAP,VOLUME,CHANGE,TOPLIST_RANK&toplist_quote_asset=USD",
		cfg.ApiURL, page, pagination)
	body, err := fetchAPIData(cfg.ApiKey, url)
	if err != nil {
		return nil, err
	}

	err = unmarshallData(body, &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %w", err)
	}

	result := make(map[int]types.AssetInfo)
	for i, v := range res.Data.List {
		v.CreatedOnFormatted = time.Unix(v.CreatedOn, 0).Format(time.RFC3339)
		v.LunchDateFormatted = time.Unix(v.LunchDate, 0).Format(time.RFC3339)
		v.PriceLastUpdatedFormatted = time.Unix(v.PriceLastUpdated, 0).Format(time.RFC3339)
		result[i] = v
	}

	return result, nil
}

func fetchAPIData(apiKey, url string) ([]byte, error) {
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

	return body, nil
}

func unmarshallData(body []byte, dest interface{}) error {
	var apiResp types.APIResponseWithError
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Err != nil && apiResp.Err.Message != "" {
		return fmt.Errorf("API error: %s", apiResp.Err.Message)
	}

	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func unmarshallCacheData(body []byte, cache map[string]types.SummaryCryptoList) error {
	var summary types.SummaryCryptoList
	err := json.Unmarshal(body, &summary)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	for _, item := range summary.Data.List {
		cache[item.Name] = summary
	}

	return nil
}
