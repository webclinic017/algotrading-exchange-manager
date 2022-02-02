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
	envOk, dbOk, kiteOk, traderOk bool
	apiKey, accToken              string
	wdg, closeTicker, initTicker  *cron.Cron
	symbolFutStr, symbolMcxFutStr string
	Tokens                        []uint32
)

func main() {

	// timeZone, _ := time.LoadLocation("Asia/Calcutta")

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

	// start watchdog to recover from connections issues
	wdg = cron.New()
	wdg.AddFunc("@every 30s", exMgrWdg)
	wdg.Start()

	// everyday scheduled start At 09:00:00 Mon-Fri
	initTicker = cron.New()
	initTicker.AddFunc("0 0 9 * * 1-5", startMainSession)
	initTicker.Start()

	// everyday scheduled stop At 16:00:00 Mon-Fri
	closeTicker = cron.New()
	closeTicker.AddFunc("0 0 16 * * 1-5", stopMainSession)
	closeTicker.Start()

	select {}

}

func startMainSession() {

	srv.InfoLogger.Print(
		"\n\n\t-------------- START ---------------",
		"\n\t------------------------------------\n\n")

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
				go db.StoreTickInDb() // TODO:check if channel open then spawn
				go trademgr.Trader()  // TODO: what condition to apply?

				// start watchdog to recover from connections issues
			}
		}
	}
	status()
}

func stopMainSession() {

	kiteOk = kite.CloseTicker()
	traderOk = trademgr.StopTrader() // Trader will terminate after closing the trades
	// DB shall close with close on channel itself - TODO: auto close logic for DB
}

func checkAPIs() {
	srv.InfoLogger.Print(
		"\n\n\t-----------------------------",
		"------------------------------------ Check API's \n\n")

	envOk = srv.LoadEnvVariables()
	dbOk = db.DbInit()
	kiteOk, apiKey, accToken = kite.LoginKite()
	status()
	db.CloseDb()
}

func exMgrWdg() {

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
