package trademgr

import (
	"algo-ex-mgr/app/apiclient"
	"algo-ex-mgr/app/appdata"
	"time"
)

func tradeEnterSignalCheck(symbol string, tradeStrategies appdata.Strategies, tr *appdata.TradeSignal) bool {

	result := false
	curTime := time.Now()

	if tradeStrategies.Trigger_time.Hour() == 0 {
		result = apiclient.SignalAnalyzer(tr, "-entry")

	} else if curTime.Hour() == tradeStrategies.Trigger_time.Hour() {
		if curTime.Minute() == tradeStrategies.Trigger_time.Minute() { // trigger time reached

			result = apiclient.SignalAnalyzer(tr, "-entry")
		}
	}

	if result {
		return true
	}

	// srv.TradesLogger.Println(" ⏳", "[", ttime, "-", sID, "-", symbol, "]")

	return false
}

func tradeExitSignalCheck(symbol string, tradeStrategies appdata.Strategies, tr *appdata.TradeSignal) bool {

	result := false
	curTime := time.Now()

	if tradeStrategies.Trigger_time.Hour() == 0 {
		result = apiclient.SignalAnalyzer(tr, "-exit")

	} else if curTime.Hour() == tradeStrategies.Trigger_time.Hour() {
		if curTime.Minute() == tradeStrategies.Trigger_time.Minute() { // trigger time reached

			result = apiclient.SignalAnalyzer(tr, "-exit")
		}
	}

	if result {
		return true
	}

	// srv.TradesLogger.Println(" ⏳", "[", ttime, "-", sID, "-", symbol, "]")

	return false
}
