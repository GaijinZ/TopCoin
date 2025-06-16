package router

import (
	"net/http"

	"topcoint/handler"
)

func Router(cryptoCurrency handler.CryptoCurrencier) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/info", cryptoCurrency.GetCryptoCurrencies)

	return router
}
