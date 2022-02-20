package main

import (
	"goTicker/app/db"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"goTicker/app/trademgr"
	"time"

	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk, traderOk bool = false, false, false, false
	apiKey, accToken              string
	wdg, sessionCron              *cron.Cron
	symbolFutStr, symbolMcxFutStr string
	Tokens                        []uint32
)

func main() {

	srv.CheckFiles()
	srv.Init()

	// testTickerData()
	// testDbFunction()

	now := time.Now()
	if (now.Hour() >= 9) && (now.Hour() < 16) &&
		(now.Weekday() > 0) && (now.Weekday() < 6) {
		startMainSession() // Check if conections are okay
	} else {
		checkAPIs() // Check if conections are okay
	}

	// everyday scheduled #[start 09:00:00] #[stop 16:00:00] (Mon-Fri)
	sessionCron = cron.New()
	sessionCron.AddFunc("0 0 9 * * 1-5", startMainSession)
	sessionCron.AddFunc("0 0 16 * * 1-5", stopMainSession)
	sessionCron.Start()

	select {}

}

func startMainSession() {

	srv.InfoLogger.Print(
		"\n\n\t-------------- START ---------------",
		"\n\t------------------------------------\n\n")

	// start watchdog to recover from connections issues
	wdg = cron.New()
	wdg.AddFunc("@every 60s", exMgrWdg)
	wdg.Start()

	envOk = srv.LoadEnvVariables()

	if envOk {

		dbOk = db.DbInit()

		if dbOk {

			kite.Tokens, kite.InsNamesMap, symbolFutStr, symbolMcxFutStr = kite.GetSymbols()
			db.StoreSymbolsInDb(symbolFutStr, symbolMcxFutStr)

			// Kite login
			kiteOk, apiKey, accToken = kite.LoginKite()

			// Start Ticker and Trader
			if kiteOk {
				kite.TickerInitialize(apiKey, accToken)

				go trademgr.Trader() // TODO: what condition to apply?
				go db.StoreTickInDb()
				// start watchdog to recover from connections issues
			}
		}
	}
	status()
}

func stopMainSession() {

	kiteOk = kite.CloseTicker() // DB will close if channel gets closed
	trademgr.StopTrader()       // Trader will terminate after closing the trades
	wdg.Stop()
}

func checkAPIs() {
	srv.InfoLogger.Print(
		"\n\n\t-----------------------------",
		"------------------------------------ Check API's \n\n")

	envOk = srv.LoadEnvVariables()
	dbOk = db.DbInit()

	trademgr.Trader()

	kiteOk, apiKey, accToken = kite.LoginKite()
	status()
	db.CloseDb()
}

func exMgrWdg() {

	// db Reconnection on error
	if (db.ErrCnt > 100) || (kite.TickerCnt < 100) {
		srv.ErrorLogger.Print("\n\n\tDB/Ticker Error, Restarting...\n\n")
		kite.CloseTicker() // close channel and DB store task
		startMainSession() // login kite, start ch & db task
	}
	db.ErrCnt = 0
	kite.TickerCnt = 0

}

func status() {
	srv.InfoLogger.Print(
		"\n\n\t--------------STATUS---------------",
		"\n\t|    ", time.Now().Format("Monday, Jan-02 3:4 PM"), "       |",
		"\n\t-----------------------------------",
		"\n\t| Environment variables set: (", envOk, ") |",
		"\n\t| DB Connected: (", dbOk, ")              |",
		"\n\t| Kite Login Succesfull: (", kiteOk, ")     |",
		"\n\t| Trader Running: (", traderOk, ")           |",
		"\n\t-----------------------------------\n\n",
	)
}

func testDbFunction() {

	kite.ChTick = make(chan kite.TickData, 1000)

	_ = srv.LoadEnvVariables()
	_ = db.DbInit()
	go db.StoreTickInDb()
	go kite.TestTicker()
	println("Testing Done")

	select {}
}

func testTickerData() {

	startMainSession()
	time.Sleep(time.Second * 20)
	stopMainSession()
	time.Sleep(time.Second * 20)
	startMainSession()
	time.Sleep(time.Second * 20)
	stopMainSession()
	println("Testing Done")

	select {}
}
