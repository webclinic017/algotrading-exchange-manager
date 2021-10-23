package main

import (
	"github.com/goTicker/kite"
)

func main() {

	// Do login and get access token
	apiKey, accToken := kite.LoginKite()

	// Initate ticker
	if accToken != "" {
		kite.TickerInitialize(apiKey, accToken)
	} else {
		println("No token generated, fail to start Ticker")
	}
}
