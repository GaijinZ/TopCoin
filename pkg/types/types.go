package types

import (
	"fmt"
)

type AssetGraph map[string][]string

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
	Symbol                    string  `json:"SYMBOL"`
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

type APIError struct {
	Message string `json:"message"`
}

type APIResponseWithError struct {
	Data map[string]AssetInfo `json:"Data"`
	Err  *APIError            `json:"Err,omitempty"`
}

type ClientMessage struct {
	Symbol     string `json:"symbol"`
	Page       int    `json:"page"`
	Pagination int    `json:"pagination"`
	Reply      chan CurrencyInfoResponse
}

func (c *ClientMessage) Validate() error {
	if c.Pagination >= 50 || c.Pagination < 10 {
		return fmt.Errorf("pagination must be less then 50 and at least 10")
	}

	if c.Page < 1 {
		return fmt.Errorf("page must be a number greater than 0")
	}

	if c.Page < 1 {
		return fmt.Errorf("page must be a number")
	}

	return nil
}
