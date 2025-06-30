package types

type CurrencyInfoRequest struct {
	Symbol     string `json:"symbol"`
	Pagination string `json:"pagination"`
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
