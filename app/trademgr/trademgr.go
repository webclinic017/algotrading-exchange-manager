// trademgr - executes and manages trades.
// Read strategies from db, spwans threads for each strategy.
// Remains active till the trade is closed
package trademgr

import (
	"algo-ex-mgr/app/apiclient"
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/srv"
	"strings"
	"sync"
	"time"
)

// tradeStrategies - list of all strategies to be executed. Read once from db at start of day
const (
	tradeOperatorSleepTime = time.Second * 2
)

var (
	TerminateTradeMgr bool = false
)

// Scan DB for all strategies with strategy_en = 1. Each funtion is executed in a separate thread and remains active till the trade is complete.
// TODO: recovery logic for server restarts
func StartTrader(daystart bool) {

	var wgTrademgr sync.WaitGroup

	srv.TradesLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Trade Manager",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	// --------------------------------- Read trading strategies from dB
	tradeUserStrategies := db.ReadUserStrategiesFromDb()

	// --------------------------------- Read if trades already in progress
	trSig := db.ReadAllOrderBookFromDb("!=", "Completed")
	for eachSymbol := range trSig {
		for eachStrategy := range tradeUserStrategies {
			if trSig[eachSymbol].Strategy == tradeUserStrategies[eachStrategy].Strategy {

				wgTrademgr.Add(1)
				go operateSymbol("nil", tradeUserStrategies[eachStrategy], trSig[eachSymbol].Id, wgTrademgr)
			}
		}
	}

	// --------------------------------- Setup operators for each symbol in every strategy
	if daystart {
		for eachStrategy := range tradeUserStrategies {

			if checkTriggerDays(tradeUserStrategies[eachStrategy]) {
				// check if the current day is a trading day.

				// Read symbols within each strategy
				tradeSymbols := strings.Split(tradeUserStrategies[eachStrategy].Instruments, ",")

				for eachSymbol := range tradeSymbols {
					wgTrademgr.Add(1)
					go operateSymbol(tradeSymbols[eachSymbol], tradeUserStrategies[eachStrategy], 0, wgTrademgr)
				}
			}
		}
	}
	// --------------------------------- Await till all trades are completed
	wgTrademgr.Wait()
}

// to stop trademanager and exit all positions
func StopTrader() {
	srv.TradesLogger.Println("(Terminating Trader) - Signal received")

}

// TODO: master exit condition & EoD termniation

// symbolTradeManager
func operateSymbol(tradeSymbol string, tradeUserStrategies appdata.UserStrategies_S, trId uint16, wgTrademgr sync.WaitGroup) {
	defer wgTrademgr.Done()

	var tr appdata.OrderBook_S
	var result bool

	if trId == 0 {
		tr.Status = "Initiate"
	} else {
		tr.Id = trId
		tr.Status = "Resume"
	}

tradingloop:
	for {
		switch tr.Status {

		// ------------------------------------------------------------------------ New symbol being registered for trade
		case "Initiate":
			tr.Date = time.Now()
			tr.Strategy = tradeUserStrategies.Strategy
			tr.Instr = tradeSymbol
			tr.Status = "AwaitSignal"
			// tr.Order_info = "{}"
			tr.Post_analysis = "{}"
			tr.Id = db.StoreOrderBookInDb(tr)

		// ------------------------------------------------------------------------ Resume previously registered symbol
		case "Resume":
			loadValues(&tr)

		// ------------------------------------------------------------------------ trade entry check (Scan Signals)
		case "AwaitSignal":
			if tradeEnterSignalCheck(tradeSymbol, tradeUserStrategies, &tr) {
				tr.Status = "PlaceOrders"
				db.StoreOrderBookInDb(tr)
			}

		// ------------------------------------------------------------------------ enter trade (order)
		case "PlaceOrders":
			if tr.Dir != "" { // on valid signal
				if tradeEnter(&tr, tradeUserStrategies) {
					tr.Status = "PlaceOrdersPending"
					db.StoreOrderBookInDb(tr)
				}
			}

			// ------------------------------------------------------------------------ enter trade (order)
		case "PlaceOrdersPending":
			if pendingOrder(&tr, tradeUserStrategies) {
				tr.Status = "TradeMonitoring"
			}
			db.StoreOrderBookInDb(tr) // store orderbook, may be partially executed

			// Todo: Add exit condition for retries

		// ------------------------------------------------------------------------ monitor trade exits
		case "TradeMonitoring":
			if apiclient.SignalAnalyzer(&tr, "-exit") {
				tr.Status = "ExitTrade"
				db.StoreOrderBookInDb(tr)
			}

		// ------------------------------------------------------------------------ squareoff trade
		case "ExitTrade":
			if result {
				tr.Status = "TradeCompleted"
				db.StoreOrderBookInDb(tr)
			}

		// ------------------------------------------------------------------------ complete housekeeping
		case "TradeCompleted":
			if result {
				db.StoreOrderBookInDb(tr)
				break tradingloop
			}

		// --------------------------------------------------------------- Terminate trade if any other status
		default:
			db.StoreOrderBookInDb(tr)
			break tradingloop
		}

		time.Sleep(tradeOperatorSleepTime)
		loadValues(&tr)
		if TerminateTradeMgr {
			tr.Status = "Terminate"
		}
		// TODO: check if exit is requested
	}
}

// Check if the current day is a trading day. Valid syntax "Monday,Tuesday,Wednesday,Thursday,Friday". For day selection to trade - Every day must be explicitly listed in dB.
func checkTriggerDays(tradeUserStrategies appdata.UserStrategies_S) bool {

	triggerdays := strings.Split(tradeUserStrategies.Trigger_days, ",")
	currentday := time.Now().Weekday().String()

	for each := range triggerdays {
		if triggerdays[each] == currentday {
			srv.TradesLogger.Println(tradeUserStrategies.Strategy, " : Trade signal registered")
			return true
		}
	}
	srv.TradesLogger.Println(tradeUserStrategies.Strategy, " : Trade signal skipped due to no valid day trigger present")
	return false
}

func loadValues(tr *appdata.OrderBook_S) {
	status, trtemp := db.ReadOrderBookFromDb(tr.Id)
	if status {
		// TODO: check if all values are loaded
		tr.Id = trtemp.Id
		tr.Date = trtemp.Date
		tr.Strategy = trtemp.Strategy
		tr.Instr = trtemp.Instr
		tr.Status = trtemp.Status
		// tr.Order_trades_entry = trtemp.Order_trades_entry
		// tr.Order_trades_exit = trtemp.Order_trades_exit
		tr.Post_analysis = trtemp.Post_analysis
	}
}

func tradeEnterSignalCheck(symbol string, tradeUserStrategies appdata.UserStrategies_S, tr *appdata.OrderBook_S) bool {

	if tradeUserStrategies.Trigger_time.Hour() == 0 {
		return apiclient.SignalAnalyzer(tr, "-entr")

	} else if time.Now().Hour() == tradeUserStrategies.Trigger_time.Hour() {
		if time.Now().Minute() == tradeUserStrategies.Trigger_time.Minute() { // trigger time reached

			return apiclient.SignalAnalyzer(tr, "-entr")
		}
	}
	return false
}
