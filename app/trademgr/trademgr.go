// trademgr - executes and manages trades.
// Read strategies from db, spwans threads for each strategy.
// Remains active till the trade is closed
package trademgr

import (
	"fmt"
	"goTicker/app/apiclient"
	"goTicker/app/data"
	"goTicker/app/db"
	"goTicker/app/srv"
	"os"
	"strings"
	"sync"
	"time"
)

// tradeStrategies - list of all strategies to be executed. Read once from db at start of day
var (
	tradeStrategies        []*data.Strategies
	terminateTradeOperator bool = false
)

// Scan DB for all strategies with strategy_en = 1. Each funtion is executed in a separate thread and remains active till the trade is complete.
// TODO: recovery logic for server restarts
func Trader() {

	var wgTrademgr sync.WaitGroup
	terminateTradeOperator = false

	srv.InitTradeLogger()

	// 1. Read trading strategies from dB
	tradeStrategies = db.ReadStrategiesFromDb()

	// 2. Setup time intervals for each strategy (loop for each)
	for each := range tradeStrategies {
		wgTrademgr.Add(1)
		go tradeOperator(tradeStrategies[each], &wgTrademgr)
	}

	// 3. wait till all trades are completed
	wgTrademgr.Wait()
	os.Exit(0)

}

func StopTrader() {
	fmt.Println("Terminating Trader - Signal received")
	terminateTradeOperator = true
}

// this is thread for each strategy
func tradeOperator(tradeStrategies *data.Strategies, wg *sync.WaitGroup) {
	defer wg.Done()

	srv.TradesLogger.Println("\n\ntradeOperator ", tradeStrategies)

	if checkTriggerDays(tradeStrategies) {
		// 1. wait for trigger time and invoke api (blocking call)
		awaitSignal(tradeStrategies)

		// 2. read db for valid signal

		// 3. on signal, execute trade (blocking call)

		// 4. on trade completion, update db

		// 5. montitor trade positions (blocking call)

		// 6. check exit conditions (blocking call)

		// 7. on signal, exit trade	(blocking call)

		// 8. on exit, update db
	}
}

// Check if the current day is a trading day. Valid syntax "Monday,Tuesday,Wednesday,Thursday,Friday". For day selection to trade - Every day must be explicitly listed in dB.
func checkTriggerDays(tradeStrategies *data.Strategies) bool {

	triggerdays := strings.Split(tradeStrategies.P_trigger_days, ",")
	currentday := time.Now().Weekday().String()

	for each := range triggerdays {
		if triggerdays[each] == currentday {
			srv.TradesLogger.Println(tradeStrategies.Strategy_id, " : Trade signal registered")
			return true
		}
	}
	srv.TradesLogger.Println(tradeStrategies.Strategy_id, " : Trade signal skipped due to no valid day trigger present")
	return false
}

// Wait till the current time is greater than the trigger time.
// TODO: master exit condition & EoD termniation
func awaitSignal(tradeStrategies *data.Strategies) {

	tradeSymbols := strings.Split(tradeStrategies.P_trade_symbols, ",")

	if tradeStrategies.P_trigger_time.Hour() == 0 {

		for {
			for each := range tradeSymbols {

				if tradeSymbols[each] != "" {
					fmt.Println("(continious) Invoking API for ", tradeStrategies.Strategy_id, "symbol : ", tradeSymbols[each])
					res, sigData := apiclient.ExecuteSingleSymbolApi(tradeStrategies.Strategy_id, tradeSymbols[each], "2022-02-09")
					if res {
						fmt.Println("Trade Signal found for ", tradeStrategies.Strategy_id, "symbol : ", tradeSymbols[each])
						db.StoreTradeSignalInDb(sigData)
						tradeSymbols[each] = "" // remove from furher scan
					}
				}
			}
			// termination requested
			if terminateTradeOperator {
				return
			}
			time.Sleep(time.Second * 10)

		}

	} else {
		// for specific time of day
		for {
			curTime := time.Now()
			triggerTime := tradeStrategies.P_trigger_time
			// fmt.Println(triggerTime, " : ", curTime)

			if curTime.Hour() == triggerTime.Hour() {
				if curTime.Minute() == triggerTime.Minute() {

					for each := range tradeSymbols {
						_, _ = apiclient.ExecuteSingleSymbolApi(tradeStrategies.Strategy_id, tradeSymbols[each], "2022-02-09")
						fmt.Println("Invoking API for ", tradeStrategies.Strategy_id, "symbol : ", tradeSymbols[each])
					}
					return
				}
			}

			// termination requested
			if terminateTradeOperator {
				return
			}

			time.Sleep(1 * time.Second * 10)
			fmt.Println("sleeping ", tradeStrategies.Strategy_id)
		}
	}
}
