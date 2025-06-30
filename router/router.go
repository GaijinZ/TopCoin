package router

import (
	"net/http"

	"topcoint/handler"
)

func Router(currencyInfo handler.Currencier) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/coindesk", currencyInfo.GetCurrencyInfo)

	fs := http.FileServer(http.Dir("public"))
	router.Handle("/", fs)

	return router
}
