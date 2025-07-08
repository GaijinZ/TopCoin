package router

import (
	"net/http"

	"topcoint/pkg/handler"
)

func Router(currencyInfo handler.Currencier) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/home", currencyInfo.Home)
	router.HandleFunc("/coindesk", currencyInfo.GetCurrencyInfo)

	router.Handle("/result.html", http.FileServer(http.Dir("public")))

	return router
}
