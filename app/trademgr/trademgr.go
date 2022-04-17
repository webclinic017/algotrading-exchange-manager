// trademgr - executes and manages trades.
// Read strategies from db, spwans threads for each strategy.
// Remains active till the trade is closed
package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/srv"
	"os"
	"strings"
	"sync"
	"time"
)

// tradeStrategies - list of all strategies to be executed. Read once from db at start of day
var (
	tradeStrategies        []*appdata.Strategies
	terminateTradeOperator bool = false
)

const (
	tradeOperatorSleepTime = time.Second * 10
)

// Scan DB for all strategies with strategy_en = 1. Each funtion is executed in a separate thread and remains active till the trade is complete.
// TODO: recovery logic for server restarts
func StartTrader() {

	var wgTrademgr sync.WaitGroup
	terminateTradeOperator = false

	srv.TradesLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Trade Manager",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	// 1. Read trading strategies from dB
	tradeStrategies = db.ReadStrategiesFromDb()

	// 2. Setup operators for each symbol in every strategy
	for eachStrategy := range tradeStrategies {

		if checkTriggerDays(tradeStrategies[eachStrategy]) { // check if the current day is a trading day.

			// Read symbols within each strategy
			tradeSymbols := strings.Split(tradeStrategies[eachStrategy].Instruments, ",")

			for eachSymbol := range tradeSymbols {
				// Check if continous OR time trigerred strategy
				wgTrademgr.Add(1)
				go operateSymbol(tradeSymbols[eachSymbol], *tradeStrategies[eachStrategy], wgTrademgr)
			}
		}
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

// TODO: master exit condition & EoD termniation

// symbolTradeManager
func operateSymbol(tradeSymbol string, tradeStrategies appdata.Strategies, wgTrademgr sync.WaitGroup) {
	defer wgTrademgr.Done()

	var orderBookId uint16
	var tr appdata.TradeSignal
	var result bool
	tr.Id = 0 // create entry in db
	tr.Date = time.Now()
	tr.Strategy = tradeStrategies.Strategy
	tr.Instr = tradeSymbol
	tr.Status = "AwaitSignal"
	tr.Order_trades_entry = "{}"
	tr.Order_trades_exit = "{}"
	tr.Order_simulation = "{}"
	tr.Post_analysis = "{}"

	orderBookId = db.StoreTradeSignalInDb(tr)
	tr.Id = orderBookId
	if orderBookId == 0 {
		srv.TradesLogger.Println("EXIT: Could not register for signal/symbol orderBookId: ", orderBookId)
		// RULE: if orderBookId is 0, then the strategy-symbol combination will not be auto traded
		return
	}

	for {
		switch tr.Status {

		// ------------------------------------------------------------------------ trade entry check (Scan Signals)
		case "AwaitSignal":
			if tradeEnterSignalCheck(tradeSymbol, tradeStrategies, &tr) {
				tr.Status = "Trigerred"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ enter trade (order)
		case "Trigerred":
			if tr.Dir != "" { // on valid signal
				result = tradeEnter(&tr, tradeStrategies)
				tr.Status = "TradeMonitoring"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ monitor trade exits
		case "TradeMonitoring":
			if tradeExitSignalCheck(tradeSymbol, tradeStrategies, &tr) {
				tr.Status = "ExitTrade"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ exit trade
		case "ExitTrade":
			if result {
				tr.Status = "TradeCompleted"
				db.StoreTradeSignalInDb(tr)
			}
		case "TradeCompleted":
			if result {
				db.StoreTradeSignalInDb(tr)
				break
			}

		default: // Terminate trade if any other status
			db.StoreTradeSignalInDb(tr)
			break

		}

		time.Sleep(tradeOperatorSleepTime)
		// read db and sync again
		// tr = db.ReadTradeSignalFromDb(orderBookId)
	}
}

// Check if the current day is a trading day. Valid syntax "Monday,Tuesday,Wednesday,Thursday,Friday". For day selection to trade - Every day must be explicitly listed in dB.
func checkTriggerDays(tradeStrategies *appdata.Strategies) bool {

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
