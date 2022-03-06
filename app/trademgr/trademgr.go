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

const (
	tradeOperatorSleepTime = time.Second * 10
	CONTINOUS_SCAN         = true
	TIME_TRIGGERED_SCAN    = false
)

// Scan DB for all strategies with strategy_en = 1. Each funtion is executed in a separate thread and remains active till the trade is complete.
// TODO: recovery logic for server restarts
func StartTrader() {

	var wgTrademgr sync.WaitGroup
	terminateTradeOperator = false

	srv.InitTradeLogger()
	srv.TradesLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Trade Manager",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

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
		tradeSymbols := strings.Split(tradeStrategies.Instruments, ",")

		for each := range tradeSymbols {

			// Check if continous OR time trigerred strategy
			if tradeStrategies.Trigger_time.Hour() == 0 {
				wgTrademgr.Add(1)
				go signalScan(CONTINOUS_SCAN, tradeSymbols[each], tradeStrategies, wgTrademgr)
			} else {
				wgTrademgr.Add(1)
				go signalScan(TIME_TRIGGERED_SCAN, tradeSymbols[each], tradeStrategies, wgTrademgr)
			}
		}
	}

	// 	// 1. wait for trigger time and invoke api (blocking call)
	// 	// 2. read db for valid signal
	// 	// 3. on signal, execute trade (blocking call)
	// 	// 4. on trade completion, update db
	// 	// 5. montitor trade positions (blocking call)
	// 	// 6. check exit conditions (blocking call)
	// 	// 7. on signal, exit trade	(blocking call)
	// 	// 8. on exit, update db
	// }
}

// Check if the current day is a trading day. Valid syntax "Monday,Tuesday,Wednesday,Thursday,Friday". For day selection to trade - Every day must be explicitly listed in dB.
func checkTriggerDays(tradeStrategies *data.Strategies) bool {

	triggerdays := strings.Split(tradeStrategies.Trigger_days, ",")
	currentday := time.Now().Weekday().String()

	for each := range triggerdays {
		if triggerdays[each] == currentday {
			srv.TradesLogger.Println(tradeStrategies.Strategy, " : Trade signal registered")
			return true
		}
	}
	srv.TradesLogger.Println(tradeStrategies.Strategy, " : Trade signal skipped due to no valid day trigger present")
	return false
}

// TODO: master exit condition & EoD termniation

// Scan signal
func signalScan(continous bool, tradeSymbol string, tradeStrategies *data.Strategies, wgTrademgr *sync.WaitGroup) {
	defer wgTrademgr.Done()

	var orderBookId uint16

	if continous {
		orderBookId = awaitContinousScan(tradeSymbol, tradeStrategies.Strategy)
	} else {
		orderBookId = awaiTriggerTimeScan(tradeSymbol, tradeStrategies.Strategy, tradeStrategies.Trigger_time)
	}
	order := db.FetchOrderData(orderBookId)

	if order != nil {
		kite.PlaceOrder(order[0])
	}
}
