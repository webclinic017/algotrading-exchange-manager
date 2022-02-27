package trademgr

import (
	"goTicker/app/apiclient"
	"goTicker/app/db"
	"goTicker/app/srv"
	"time"
)

func awaitContinousScan(symbol string, sID string) uint16 {

	var orderBookId uint16 = 0

	for {
		srv.TradesLogger.Println("(Continious Scan) Invoking API [", sID, "-", symbol, "]")
		result, sigData := apiclient.SignalAnalyzer("false", sID, symbol, "2022-02-09")

		if result {
			srv.TradesLogger.Println("{Trade Signal found} [", sID, "-", symbol, "]")
			orderBookId = db.StoreTradeSignalInDb(sigData)
			break
		}

		if terminateTradeOperator { // termination requested
			srv.TradesLogger.Println("(Continious Scan) Termination requested [", sID, "-", symbol, "]")
			break
		}
		time.Sleep(tradeOperatorSleepTime)
	}
	return orderBookId
}

func awaiTriggerTimeScan(symbol string, sID string, triggerTime time.Time) uint16 {

	var orderBookId uint16 = 0

	for {
		curTime := time.Now()

		// fmt.Println(triggerTime, " : ", curTime)

		if curTime.Hour() == triggerTime.Hour() {
			if curTime.Minute() == triggerTime.Minute() { // trigger time reached

				srv.TradesLogger.Println("Invoking TimeTrigerred API [", sID, "-", symbol, "]")
				result, sigData := apiclient.SignalAnalyzer("false", sID, symbol, "2022-02-09")

				if result {
					srv.TradesLogger.Println("{Trade Signal found} [", sID, "-", symbol, "]")
					orderBookId = db.StoreTradeSignalInDb(sigData)
				}
				break
			}
		}

		// termination requested
		if terminateTradeOperator {
			srv.TradesLogger.Println("(TimeTrigerred) Termination requested ", sID, "symbol : ", symbol)
			return 0
		}

		time.Sleep(tradeOperatorSleepTime)
		srv.TradesLogger.Println("(Sleeping) [", sID, "-", symbol, "]")
	}

	return orderBookId
}
