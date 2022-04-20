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
	tradeStrategies := db.ReadStrategiesFromDb()

	// --------------------------------- Read if trades already in progress
	trSig := db.ReadAllTradeSignalFromDb("!=", "Completed")
	for eachSymbol := range trSig {
		for eachStrategy := range tradeStrategies {
			if trSig[eachSymbol].Strategy == tradeStrategies[eachStrategy].Strategy {

				wgTrademgr.Add(1)
				go operateSymbol("nil", tradeStrategies[eachStrategy], trSig[eachSymbol].Id, wgTrademgr)
			}
		}
	}

	// --------------------------------- Setup operators for each symbol in every strategy
	if daystart {
		for eachStrategy := range tradeStrategies {

			if checkTriggerDays(tradeStrategies[eachStrategy]) {
				// check if the current day is a trading day.

				// Read symbols within each strategy
				tradeSymbols := strings.Split(tradeStrategies[eachStrategy].Instruments, ",")

				for eachSymbol := range tradeSymbols {
					wgTrademgr.Add(1)
					go operateSymbol(tradeSymbols[eachSymbol], tradeStrategies[eachStrategy], 0, wgTrademgr)
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
func operateSymbol(tradeSymbol string, tradeStrategies appdata.Strategies, trId uint16, wgTrademgr sync.WaitGroup) {
	defer wgTrademgr.Done()

	var tr appdata.TradeSignal
	var result bool

	if trId == 0 {
		tr.Status = "Initiate"
	} else {
		tr.Status = "Resume"
	}

tradingloop:
	for {
		switch tr.Status {

		// ------------------------------------------------------------------------ New symbol being registered for trade
		case "Initiate":
			tr.Date = time.Now()
			tr.Strategy = tradeStrategies.Strategy
			tr.Instr = tradeSymbol
			tr.Status = "AwaitSignal"
			tr.Order_trade_entry = "{}"
			tr.Order_trade_exit = "{}"
			tr.Order_simulation = "{}"
			tr.Post_analysis = "{}"
			tr.Status = "AwaitSignal"
			tr.Id = db.StoreTradeSignalInDb(tr)

		// ------------------------------------------------------------------------ Resume previously registered symbol
		case "Resume":
			loadValues(&tr)
			db.StoreTradeSignalInDb(tr)

		// ------------------------------------------------------------------------ trade entry check (Scan Signals)
		case "AwaitSignal":
			if tradeEnterSignalCheck(tradeSymbol, tradeStrategies, &tr) {
				tr.Status = "PlaceOrders"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ enter trade (order)
		case "PlaceOrders":
			if tr.Dir != "" { // on valid signal
				result = tradeEnter(&tr, tradeStrategies)
				tr.Status = "TradeMonitoring"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ monitor trade exits
		case "TradeMonitoring":
			if apiclient.SignalAnalyzer(&tr, "-exit") {
				tr.Status = "ExitTrade"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ squareoff trade
		case "ExitTrade":
			if result {
				tr.Status = "TradeCompleted"
				db.StoreTradeSignalInDb(tr)
			}

		// ------------------------------------------------------------------------ complete housekeeping
		case "TradeCompleted":
			if result {
				db.StoreTradeSignalInDb(tr)
				break tradingloop
			}

		// --------------------------------------------------------------- Terminate trade if any other status
		default:
			db.StoreTradeSignalInDb(tr)
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
func checkTriggerDays(tradeStrategies appdata.Strategies) bool {

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

func loadValues(tr *appdata.TradeSignal) {
	status, trtemp := db.ReadTradeSignalFromDb(tr.Id)
	if status {
		tr.Id = trtemp.Id
		tr.Date = trtemp.Date
		tr.Strategy = trtemp.Strategy
		tr.Instr = trtemp.Instr
		tr.Status = trtemp.Status
		tr.Order_trade_entry = trtemp.Order_trade_entry
		tr.Order_trade_exit = trtemp.Order_trade_exit
		tr.Order_simulation = trtemp.Order_simulation
		tr.Post_analysis = trtemp.Post_analysis
	}

}

func tradeEnterSignalCheck(symbol string, tradeStrategies appdata.Strategies, tr *appdata.TradeSignal) bool {

	if tradeStrategies.Trigger_time.Hour() == 0 {
		return apiclient.SignalAnalyzer(tr, "-entry")

	} else if time.Now().Hour() == tradeStrategies.Trigger_time.Hour() {
		if time.Now().Minute() == tradeStrategies.Trigger_time.Minute() { // trigger time reached

			return apiclient.SignalAnalyzer(tr, "-entry")
		}
	}
	return false
}
