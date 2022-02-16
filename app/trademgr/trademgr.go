package trademgr

import (
	"goTicker/app/data"
	"goTicker/app/db"
)

var tradeStrategies []*data.Strategies

func Trader() {

	// 1. Read trading strategies from dB
	tradeStrategies = db.ReadStrategiesFromDb()

	// 2. Setup time intervals for each strategy (loop for each)
	print(tradeStrategies[0].Strategy_id)

	// 3. Exit monitoring - every 1 minute

}

func StopTrader() bool {

	return true
}
