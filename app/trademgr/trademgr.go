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
	awaitSignalSleep = time.Second * 2
	placeOrderSleep  = time.Millisecond * 100
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

	// --------------------------------- Resume operations on restart or new day start
	trSig := db.ReadAllOrderBookFromDb("!=", "Completed")
	var s bool = false
	for eachSymbol := range trSig {
		s = false
		for eachStrategy := range tradeUserStrategies {
			if trSig[eachSymbol].Strategy == tradeUserStrategies[eachStrategy].Strategy {

				wgTrademgr.Add(1)
				srv.TradesLogger.Println(appdata.ColorSuccess, "\n\nStrategy being resumed\n", trSig[eachSymbol])
				go operateSymbol("nil", tradeUserStrategies[eachStrategy], trSig[eachSymbol].Id, wgTrademgr)
				s = true
				break
			}
		}
		if !s {
			srv.TradesLogger.Println(appdata.ColorError, "\n\nStrategy could not be resumed\n", trSig[eachSymbol])
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
	TerminateTradeMgr = true
	srv.TradesLogger.Println("(Terminating Trader) - Signal received")

}

// TODO: master exit condition & EoD termniation

// symbolTradeManager
func operateSymbol(tradeSymbol string, tradeUserStrategies appdata.UserStrategies_S, orderId uint16, wgTrademgr sync.WaitGroup) {
	defer wgTrademgr.Done()

	var order appdata.OrderBook_S
	var result bool

	if orderId == 0 {
		order.Status = "Initiate"
	} else { // Resume previously registered symbol
		order.Id = orderId
		loadValues(&order)
	}

tradingloop:
	for {
		switch order.Status {

		// ------------------------------------------------------------------------ New symbol being registered for trade
		case "Initiate":
			order.Date = time.Now()
			order.Strategy = tradeUserStrategies.Strategy
			order.Instr = tradeSymbol
			order.Status = "AwaitSignal"
			order.Info.Order_simulation = tradeUserStrategies.Parameters.Controls.TradeSimulate
			// tr.Order_info = "{}"
			order.Post_analysis = "{}"
			order.Id = db.StoreOrderBookInDb(order)

		// ------------------------------------------------------------------------ trade entry check (Scan Signals)
		case "AwaitSignal":
			if tradeEnterSignalCheck(tradeSymbol, tradeUserStrategies, &order) {
				order.Status = "PlaceOrders"
				db.StoreOrderBookInDb(order)
			}
			time.Sleep(awaitSignalSleep)

		// ------------------------------------------------------------------------ enter trade (order)
		case "PlaceOrders":
			if order.Dir != "" { // on valid signal
				if tradeEnter(&order, tradeUserStrategies) {
					order.Status = "PlaceOrdersPending"
					db.StoreOrderBookInDb(order)
				}
			}
			time.Sleep(placeOrderSleep)

			// ------------------------------------------------------------------------ enter trade (order)
		case "PlaceOrdersPending":
			if pendingOrder(&order, tradeUserStrategies) {
				order.Status = "TradeMonitoring"
			}
			db.StoreOrderBookInDb(order) // store orderbook, may be partially executed
			time.Sleep(placeOrderSleep)
			// Todo: Add exit condition for retries

		// ------------------------------------------------------------------------ monitor trade exits
		case "TradeMonitoring":
			if apiclient.SignalAnalyzer(&order, "-exit") {
				order.Status = "ExitTrade"
				db.StoreOrderBookInDb(order)
			}
			time.Sleep(awaitSignalSleep)

		// ------------------------------------------------------------------------ squareoff trade
		case "ExitTrade":
			if tradeExit(&order, tradeUserStrategies) {
				order.Status = "ExitOrdersPending"
				db.StoreOrderBookInDb(order)
			}
			time.Sleep(placeOrderSleep)

			// ------------------------------------------------------------------------ enter trade (order)
		case "ExitOrdersPending":
			if pendingOrder(&order, tradeUserStrategies) {
				order.Status = "TradeCompleted"
			}
			db.StoreOrderBookInDb(order) // store orderbook, may be partially executed
			time.Sleep(awaitSignalSleep)

			// Todo: Add exit condition for retries

		// ------------------------------------------------------------------------ complete housekeeping
		case "TradeCompleted":
			if result {
				db.StoreOrderBookInDb(order)
				break tradingloop
			}
			time.Sleep(awaitSignalSleep)

		// --------------------------------------------------------------- Terminate trade if any other status
		default:
			db.StoreOrderBookInDb(order)
			break tradingloop
		}

		loadValues(&order)
		if TerminateTradeMgr {
			order.Status = "Terminate"
		} else if order.Info.UserExitRequested {
			order.Status = "ExitTrade"
		}
	}
}

// RULE: Check if the current day is a trading day. Valid syntax "Monday,Tuesday,Wednesday,Thursday,Friday". For day selection to trade - Every day must be explicitly listed in dB.
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

func loadValues(or *appdata.OrderBook_S) {
	status, trtemp := db.ReadOrderBookFromDb(or.Id)
	if status {
		or.Id = trtemp.Id
		or.Date = trtemp.Date
		or.Instr = trtemp.Instr
		or.Strategy = trtemp.Strategy
		or.Status = trtemp.Status
		or.Dir = trtemp.Dir
		or.Exit_reason = trtemp.Exit_reason
		or.Info = trtemp.Info
		or.Targets = trtemp.Targets
		or.Orders_entr = trtemp.Orders_entr
		or.Orders_exit = trtemp.Orders_exit
		or.Post_analysis = trtemp.Post_analysis
	} else {
		or.Info.ErrorCount++
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
