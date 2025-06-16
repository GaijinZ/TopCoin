package router

import (
	"net/http"

	"topcoint/coindesk/handler"
)

func Router(currencyInfo handler.Currencier) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/coindesk", currencyInfo.GetCurrencyInfo)

	return router
}
