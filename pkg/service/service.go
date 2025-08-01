package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"topcoint/pkg/config"
	"topcoint/pkg/types"
)

type Cacher interface {
	HandleIncomingMessages(cfg config.Config, msg types.ClientMessage) (types.CurrencyInfoResponse, error)
	CacheCurrency(cfg config.Config) error
	BuildCurrencyInfoResponse(cfg config.Config, msg types.ClientMessage) (types.CurrencyInfoResponse, error)
}

type CryptoCurrencyList struct {
	Symbol string `json:"SYMBOL"`
	Name   string `json:"NAME"`
}

type SummaryCryptoList struct {
	Data struct {
		List []CryptoCurrencyList `json:"LIST"`
	} `json:"Data"`
}

func NewSummaryCryptoList() Cacher {
	return &SummaryCryptoList{
		Data: struct {
			List []CryptoCurrencyList `json:"LIST"`
		}{
			List: make([]CryptoCurrencyList, 0),
		},
	}
}

func (c *SummaryCryptoList) HandleIncomingMessages(cfg config.Config, req types.ClientMessage) (types.CurrencyInfoResponse, error) {
	response, err := c.BuildCurrencyInfoResponse(cfg, req)
	if err != nil {
		return types.CurrencyInfoResponse{}, fmt.Errorf("error building currency info response: %w", err)
	}

	return response, nil
}

func (c *SummaryCryptoList) CacheCurrency(cfg config.Config) error {
	url := fmt.Sprintf("%s/asset/v1/summary/list?asset_lookup_priority=SYMBOL", cfg.ApiURL)
	body, err := fetchAPIData(cfg.ApiKey, url)
	if err != nil {
		return fmt.Errorf("failed to fetch currency data: %w", err)
	}

	err = unmarshallCacheData(body, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

func (c *SummaryCryptoList) BuildCurrencyInfoResponse(cfg config.Config, msg types.ClientMessage) (types.CurrencyInfoResponse, error) {
	response := types.CurrencyInfoResponse{Data: make(map[string]interface{})}
	input := strings.ToUpper(msg.Symbol)

	matchedSymbols := findSuggestions(c.Data.List, input)
	perfectMatchSymbol := findPerfectSymbolMatch(c.Data.List, input)

	if perfectMatchSymbol == "" {
		response.Data["symbols"] = matchedSymbols
		return response, nil
	}

	info, err := fetchCurrencyInfo(cfg, perfectMatchSymbol)
	if err != nil {
		return response, fmt.Errorf("error fetching currency info: %w", err)
	}

	for k, v := range info {
		response.Data[k] = v
	}

	if msg.Pagination > 0 && msg.Page > 0 {
		if err := msg.Validate(); err != nil {
			return response, fmt.Errorf("invalid message: %w", err)
		}

		stats, err := fetchCurrencyStats(cfg, msg.Page, msg.Pagination)
		if err != nil {
			return response, fmt.Errorf("error fetching currency stats: %w", err)
		}

		response.Data["stats"] = stats
	}

	return response, nil
}

func findSuggestions(list []CryptoCurrencyList, input string) []map[string]string {
	suggestions := make([]map[string]string, 0)

	for _, entry := range list {
		symbol := strings.ToUpper(entry.Symbol)
		name := strings.ToUpper(entry.Name)
		if strings.HasPrefix(symbol, input) || strings.HasPrefix(name, input) {
			suggestions = append(suggestions, map[string]string{
				"symbol": entry.Symbol,
				"name":   entry.Name,
			})
		}
	}

	return suggestions
}

func findPerfectSymbolMatch(list []CryptoCurrencyList, input string) string {
	for _, entry := range list {
		if strings.ToUpper(entry.Symbol) == input {
			return entry.Symbol
		}
	}

	return ""
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

func unmarshallCacheData(body []byte, cache *SummaryCryptoList) error {
	err := json.Unmarshal(body, &cache)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
