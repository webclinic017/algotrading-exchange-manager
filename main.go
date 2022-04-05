package main

import (
	"algo-ex-mgr/app/apiclient"
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"os"
	"time"

	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk, traderOk bool = false, false, false, false
	wdg, sessionCron              *cron.Cron
)

func main() {

	srv.Init()
	now := time.Now()
	if (now.Hour() >= 9) && (now.Hour() < 16) &&
		(now.Weekday() > 0) && (now.Weekday() < 6) {
		startMainSession() // start if App invoked in trade time
	} else {
		checkAPIs() // Check if conections are okay
	}

	// everyday scheduled #[start 09:00:00] #[stop 16:00:00] (Mon-Fri)
	sessionCron = cron.New()
	sessionCron.AddFunc("0 0 9 * * 1-5", startMainSession)
	sessionCron.AddFunc("0 0 16 * * 1-5", stopMainSession)
	sessionCron.AddFunc("0 0 8 * * 1-5", preTradeOps)
	sessionCron.AddFunc("0 5 16 * * 1-5", postTradeOps)
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

	envOk = srv.LoadEnvVariables("./userSettings.env")
	if envOk {

		dbOk = db.DbInit()
		if dbOk {
			// Kite login
			kiteOk = kite.Init()
			// Start Ticker and Trader
			if kiteOk {
				kite.TickerInitialize(appdata.Env["ZERODHA_API_KEY"], os.Getenv("kiteaccessToken"))

				// go trademgr.StartTrader() // TODO: what condition to apply?
				db.InitTickStorage()
				// start watchdog to recover from connections issues
			}
		}
	}
	status()
}

func stopMainSession() {

	kiteOk = kite.CloseTicker() // DB will close if channel gets closed
	// trademgr.StopTrader()       // Trader will terminate after closing the trades
	wdg.Stop()
}

func checkAPIs() {
	srv.InfoLogger.Print(
		"\n\n\t-----------------------------",
		"------------------------------------ Check API's --- MARKET OFF-TIME\n\n")

	// startMainSession()
	// go trademgr.StartTrader()
	// trademgr.StopTrader()
	// time.Sleep(time.Second * 5)
	// stopMainSession()
	// os.Exit(0)

	envOk = srv.LoadEnvVariables("./userSettings.env")
	dbOk = db.DbInit()
	kiteOk = kite.Init()
	status()
	if dbOk {
		db.CloseDb()
	}

}

func exMgrWdg() {

	// db Reconnection on error
	if (db.ErrCnt > 100) || (kite.TickerCnt < 100) {
		srv.ErrorLogger.Print("\n\n\tWatchdogMgr - | db.ErrCnt:", db.ErrCnt, "\tkite.TickerCnt:", kite.TickerCnt, "\n\n")
		// srv.ErrorLogger.Print("\n\n\tWatchdog - DB/Ticker Error, Restarting...\n\n")
		// kite.CloseTicker() // close channel and DB store task
		// startMainSession() // login kite, start ch & db task
		time.Sleep(time.Minute * 1) // wait to establish connections
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

	// appdata.ChTick = make(chan appdata.TickData, 1000)

	// _ = srv.LoadEnvVariables("app/zfiles/config/userSettings.env")
	// _ = db.DbInit()
	// go db.StoreTickInDb()
	// go kite.TestTicker()
	// println("Testing Done")

	// select {}
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

func preTradeOps() {
	srv.InfoLogger.Println("preTradeOps Started")
	if !apiclient.Services("instruments", time.Now()) {
		srv.ErrorLogger.Println("FAILED - preTradeOps/instruments")
	}

}

func postTradeOps() {
	srv.InfoLogger.Println("postTradeOps Started")
	if !apiclient.Services("candle1min-converter", time.Now()) {
		srv.ErrorLogger.Println("FAILED - postTradeOps/candle1min-converter")
	}
}
