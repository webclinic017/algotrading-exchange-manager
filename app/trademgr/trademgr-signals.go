package trademgr

import (
	"algo-ex-mgr/app/apiclient"
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"strconv"
	"time"
)

func signalAwaitContinous(symbol string, sID string, tr *appdata.TradeSignal) string {

	for {
		time.Sleep(tradeOperatorSleepTime)

		srv.TradesLogger.Println(" ▶ (Continious Scan) Invoking API [", sID, "-", symbol, "]")
		result, sigData := apiclient.SignalAnalyzer("false", sID, symbol, "2022-02-09")

		if result {
			tr.Status = "Continous - Trade Signal found"
			srv.TradesLogger.Println(" ⏪ {Trade Signal found} [", sID, "-", symbol, "]")
			return sigData
		}

		if terminateTradeOperator { // termination requested
			tr.Status = "Continous - Termination requested"
			srv.TradesLogger.Println(" ❎ (Continious Scan) Termination requested [", sID, "-", symbol, "]")
			return ""
		}

	}
}

func signalAwaitTimeTrigerred(symbol string, sID string, triggerTime time.Time, tr *appdata.TradeSignal) string {

	ttime := strconv.Itoa(int(triggerTime.Hour())) + ":" + strconv.Itoa(int(triggerTime.Minute()))

	for {
		curTime := time.Now()

		if curTime.Hour() == triggerTime.Hour() {
			if curTime.Minute() == triggerTime.Minute() { // trigger time reached

				srv.TradesLogger.Println(" ▶ Invoking TimeTrigerred API [ (", ttime, ") -", sID, "-", symbol, "]")
				result, sigData := apiclient.SignalAnalyzer("false", sID, symbol, "2022-02-09")

				if result {
					tr.Status = "TimeTrigerred - Trade Signal found"
					srv.TradesLogger.Println(" ⏪ {Trade Signal found}",
						"[",
						ttime,
						"-",
						sID,
						"-",
						symbol,
						"]")
					return sigData
				}
				break
			}
		}

		// termination requested
		if terminateTradeOperator {
			tr.Status = "TimeTrigerred - Termination requested"
			srv.TradesLogger.Println(" ❎ (TimeTrigerred) Termination requested ", sID, "symbol : ", symbol)
			return ""
		}

		time.Sleep(tradeOperatorSleepTime)
		srv.TradesLogger.Println(" ⏳",
			"[",
			ttime,
			"-",
			sID,
			"-",
			symbol,
			"]")
	}
	return ""
}
