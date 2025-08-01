package router

import (
	"net/http"
	"topcoint/pkg/types"

	"topcoint/pkg/handler"
)

func Router(requests chan<- types.ClientMessage, currencyInfo handler.Currencier) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/home", currencyInfo.Home)
	router.HandleFunc("/coindesk", currencyInfo.WSHandler(requests))

	router.Handle("/", http.FileServer(http.Dir("./public")))

	return router
}
