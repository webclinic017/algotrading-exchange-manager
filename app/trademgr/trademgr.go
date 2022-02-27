// trademgr - executes and manages trades.
// Read strategies from db, spwans threads for each strategy.
// Remains active till the trade is closed
package trademgr

import (
	"goTicker/app/data"
	"goTicker/app/db"
	"goTicker/app/kite"
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

const tradeOperatorSleepTime = time.Second * 10

// Scan DB for all strategies with strategy_en = 1. Each funtion is executed in a separate thread and remains active till the trade is complete.
// TODO: recovery logic for server restarts
func StartTrader() {

	var wgTrademgr sync.WaitGroup
	terminateTradeOperator = false

	srv.InitTradeLogger()

	// 1. Read trading strategies from dB
	tradeStrategies = db.ReadStrategiesFromDb()

	// 2. Setup operators for each strategy
	for each := range tradeStrategies {
		tradeOperator(tradeStrategies[each], &wgTrademgr) // for each strategy
	}

	// 3. wait till all trades are completed
	wgTrademgr.Wait()
	os.Exit(0)

}

// to stop trademanager and exit all positions
func StopTrader() {
	srv.TradesLogger.Println("(Terminating Trader) - Signal received")
	terminateTradeOperator = true
}

// Scans for all strategies and spawn thread for each symbol in that strategy
func tradeOperator(tradeStrategies *data.Strategies, wgTrademgr *sync.WaitGroup) {

	srv.TradesLogger.Println("\n(TradeOperator Setup) ", tradeStrategies)

	if checkTriggerDays(tradeStrategies) { // check if the current day is a trading day.

		// Read symbols within each strategy
		tradeSymbols := strings.Split(tradeStrategies.P_trade_symbols, ",")

		for each := range tradeSymbols {

			// Check if continous OR time trigerred strategy
			if tradeStrategies.P_trigger_time.Hour() == 0 {
				wgTrademgr.Add(1)
				go toContinous(tradeSymbols[each], tradeStrategies, wgTrademgr)
			} else {
				wgTrademgr.Add(1)
				go toTimeTrigerred(tradeSymbols[each], tradeStrategies, wgTrademgr)
			}
		}

		// 1. wait for trigger time and invoke api (blocking call)
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

// TODO: master exit condition & EoD termniation

// Continous scan strategy
func toContinous(tradeSymbol string, tradeStrategies *data.Strategies, wgTrademgr *sync.WaitGroup) {
	defer wgTrademgr.Done()

	orderBookId := awaitContinousScan(tradeSymbol, tradeStrategies.Strategy_id)
	// fetchRecord(orderBookId)
	kite.ExecuteTrade(orderBookId)

}

// Strategy invoked at the time of trigger.
func toTimeTrigerred(tradeSymbol string, tradeStrategies *data.Strategies, wgTrademgr *sync.WaitGroup) {
	defer wgTrademgr.Done()

	orderBookId := awaiTriggerTimeScan(tradeSymbol, tradeStrategies.Strategy_id, tradeStrategies.P_trigger_time)
	// fetchRecord(orderBookId)
	kite.ExecuteTrade(orderBookId)

}
