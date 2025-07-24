package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type CurrencyInfoResponse struct {
	Data map[string]interface{} `json:"Data"`
}

type CurrencyInfoAPIResponse struct {
	Data map[string]AssetInfo `json:"Data"`
}

type CurrencyStatsAPIResponse struct {
	Data CurrencyAPIStatsData `json:"Data"`
}

type CurrencyAPIStatsData struct {
	List []AssetInfo `json:"List"`
}

type AssetInfo struct {
	CreatedOn                 int64   `json:"CREATED_ON"`
	CreatedOnFormatted        string  `json:"CREATED_ON_FORMATTED,omitempty"`
	LunchDate                 int64   `json:"LAUNCH_DATE"`
	LunchDateFormatted        string  `json:"LUNCH_DATE_FORMATTED,omitempty"`
	AssetType                 string  `json:"ASSET_TYPE"`
	Name                      string  `json:"NAME"`
	Price                     float64 `json:"PRICE_USD"`
	PriceLastUpdated          int64   `json:"PRICE_USD_LAST_UPDATE_TS"`
	PriceLastUpdatedFormatted string  `json:"PRICE_USD_LAST_FORMATTED,omitempty"`
	Description               string  `json:"ASSET_DESCRIPTION"`
}

type CryptoCurrencyList struct {
	Name string `json:"NAME"`
}

type SummaryCryptoList struct {
	Data struct {
		List []CryptoCurrencyList `json:"LIST"`
	} `json:"Data"`
}

type APIError struct {
	Message string `json:"message"`
}

type APIResponseWithError struct {
	Data map[string]AssetInfo `json:"Data"`
	Err  *APIError            `json:"Err,omitempty"`
}

type ClientMessage struct {
	Action     string `json:"action"`
	Symbol     string `json:"symbol"`
	Page       int    `json:"page"`
	Pagination int    `json:"pagination"`
}

func (c *ClientMessage) Parse(body []byte) error {
	type tmpMsg struct {
		Action     string `json:"action"`
		Symbol     string `json:"symbol"`
		Page       string `json:"page"`
		Pagination string `json:"pagination"`
	}

	var tmp tmpMsg
	if err := json.Unmarshal(body, &tmp); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	page, err := strconv.Atoi(tmp.Page)
	if err != nil {
		return fmt.Errorf("invalid page: %w", err)
	}

	pagination, err := strconv.Atoi(tmp.Pagination)
	if err != nil {
		return fmt.Errorf("invalid pagination: %w", err)
	}

	c.Action = tmp.Action
	c.Symbol = tmp.Symbol
	c.Page = page
	c.Pagination = pagination

	return nil
}

func (c *ClientMessage) Validate() error {
	if c.Pagination >= 50 || c.Pagination < 10 {
		return fmt.Errorf("pagination must be less then 50 and at least 10")
	}

	if c.Page < 1 {
		return fmt.Errorf("page must be a number")
	}

	return nil
}
