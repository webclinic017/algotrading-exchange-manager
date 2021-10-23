package main

import (
	"time"

	"github.com/goTicker/dbticks"
	"github.com/goTicker/kite"
)

func main() {

	dbticks.DbInit()

	// Do login and get access token
	apiKey, accToken := kite.LoginKite()

	// Initate ticker
	if accToken != "" {
		kite.TickerInitialize(apiKey, accToken)
	} else {
		println("No token generated, fail to start Ticker")
	}
	time.Sleep(5 * time.Second)
	kite.CloseTicker()
	time.Sleep(5 * time.Second)

}
