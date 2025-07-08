package types

import (
	"fmt"
	"net/http"
	"strconv"
)

type CurrencyInfoRequest struct {
	Symbol     string `json:"symbol"`
	Pagination string `json:"pagination"`
	Page       string `json:"page"`
}

type CurrencyInfoResponse struct {
	Data map[string]interface{} `json:"Data"`
}

type CurrencyInfoAPIResponse struct {
	Data map[string]interface{} `json:"Data"`
}

type CurrencyStatsAPIResponse struct {
	Data map[string]interface{} `json:"Data"`
}

type APIError struct {
	Message string `json:"message"`
}

type APIResponseWithError struct {
	Data map[string]interface{} `json:"Data"`
	Err  *APIError              `json:"Err,omitempty"`
}

type ClientMessage struct {
	Action     string `json:"action"`
	Symbol     string `json:"symbol"`
	Page       string `json:"page"`
	Pagination string `json:"pagination"`
}

func (c *CurrencyInfoRequest) Parse(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	c.Symbol = r.FormValue("symbol")
	c.Pagination = r.FormValue("pagination")
	c.Page = r.FormValue("page")

	return nil
}

func (c *CurrencyInfoRequest) Validate() error {
	if c.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	if c.Pagination == "" {
		return fmt.Errorf("pagination is required")
	}

	if c.Page == "" {
		return fmt.Errorf("page is required")
	}

	if _, err := strconv.Atoi(c.Pagination); err != nil {
		return fmt.Errorf("pagination must be a number")
	}

	if _, err := strconv.Atoi(c.Page); err != nil {
		return fmt.Errorf("page must be a number")
	}

	if c.Pagination < "10" {
		return fmt.Errorf("pagination must be at least 10")
	}

	return nil
}
