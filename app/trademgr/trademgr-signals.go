package trademgr

import (
	"goTicker/app/apiclient"
	"goTicker/app/db"
	"goTicker/app/srv"
	"strconv"
	"time"
)

func awaitContinousScan(symbol string, sID string) uint16 {

	var orderBookId uint16 = 0

	for {
		time.Sleep(tradeOperatorSleepTime)

		srv.TradesLogger.Println(" ▶ (Continious Scan) Invoking API [", sID, "-", symbol, "]")
		result, sigData := apiclient.SignalAnalyzer("false", sID, symbol, "2022-02-09")

		if result {
			srv.TradesLogger.Println(" ⏪ {Trade Signal found} [", sID, "-", symbol, "]")
			orderBookId = db.StoreTradeSignalInDb(sigData)
			break
		}

		if terminateTradeOperator { // termination requested
			srv.TradesLogger.Println(" ❎ (Continious Scan) Termination requested [", sID, "-", symbol, "]")
			break
		}

	}
	return orderBookId
}

func awaiTriggerTimeScan(symbol string, sID string, triggerTime time.Time) uint16 {

	var orderBookId uint16 = 0
	ttime := strconv.Itoa(int(triggerTime.Hour())) + ":" + strconv.Itoa(int(triggerTime.Minute()))

	for {
		curTime := time.Now()

		if curTime.Hour() == triggerTime.Hour() {
			if curTime.Minute() == triggerTime.Minute() { // trigger time reached

				srv.TradesLogger.Println(" ▶ Invoking TimeTrigerred API [ (", ttime, ") -", sID, "-", symbol, "]")
				result, sigData := apiclient.SignalAnalyzer("false", sID, symbol, "2022-02-09")

				if result {
					srv.TradesLogger.Println(" ⏪ {Trade Signal found}",
						"[",
						ttime,
						"-",
						sID,
						"-",
						symbol,
						"]")
					orderBookId = db.StoreTradeSignalInDb(sigData)
				}
				break
			}
		}

		// termination requested
		if terminateTradeOperator {
			srv.TradesLogger.Println(" ❎ (TimeTrigerred) Termination requested ", sID, "symbol : ", symbol)
			return 0
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

	return orderBookId
}
