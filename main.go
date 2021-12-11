package main

import (
	"goTicker/app/cdlconv"
	"goTicker/app/db"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk                                     bool
	apiKey, accToken                                        string
	cdl1min, cdl3min, cdl5min, wdg, closeTicker, initTicker *cron.Cron
	Tokens                                                  []uint32
)

func main() {

	// timeZone, _ := time.LoadLocation("Asia/Calcutta")

	srv.Init()

	kite.Tokens = kite.GetSymbols()

	now := time.Now()

	if (now.Hour() >= 9) && (now.Hour() < 16) &&
		(now.Weekday() > 0) && (now.Weekday() < 6) {
		initTickerToken() // start now, when docker starts if its within trading time (9am-4pm Mon-Fri)
	} else {
		firstRunConnectionsCheck() // Check if conections are okay
	}
	// everyday scheduled start
	initTicker = cron.New()
	initTicker.AddFunc("0 0 9 * * 1-5", initTickerToken) // At 09:00:00 Mon-Fri
	initTicker.Start()

	// everyday scheduled stop
	closeTicker = cron.New()
	closeTicker.AddFunc("0 16 * * 1-5", theStop) // At 16:00:00 Mon-Fri
	closeTicker.Start()

	select {}

}

func loadEnv() bool {

	// Load .env file, if not in production
	if 0 >= len(os.Getenv("PRODUCTION")) {
		err := godotenv.Load("app/config/ENV_Settings.env")
		if err != nil {
			srv.ErrorLogger.Fatal("ENV_Settings.env file not found, Terminating!!!")
			//return false
		}
	}
	if 0 >= len(os.Getenv("TFA_AUTH")) {
		srv.ErrorLogger.Println("TFA_AUTH not set")
		return false
	}
	if 0 >= len(os.Getenv("USER_ID")) {
		srv.ErrorLogger.Println("USER_ID not set")
		return false
	}
	if 0 >= len(os.Getenv("PASSWORD")) {
		srv.ErrorLogger.Println("PASSWORD not set")
		return false
	}
	if 0 >= len(os.Getenv("API_KEY")) {
		srv.ErrorLogger.Println("API_KEY not set")
		return false
	}
	if 0 >= len(os.Getenv("API_SECRET")) {
		srv.ErrorLogger.Println("API_SECRET not set")
		return false
	}
	if 0 >= len(os.Getenv("REQUEST_TOKEN_URL")) {
		srv.ErrorLogger.Println("REQUEST_TOKEN_URL not set")
		return false
	}
	if 0 >= len(os.Getenv("EXTERNAL_DATABASE_URL")) {
		srv.ErrorLogger.Println("EXTERNAL_DATABASE_URL not set")
		return false
	}
	os.Setenv("DATABASE_URL", os.Getenv("EXTERNAL_DATABASE_URL"))

	return true
}

func printStatus(envOk, dbOk, kiteOk bool) {
	srv.InfoLogger.Printf("\n--------STATUS---------\nEnvironment files found: %t\nKite Login Succesfull: %t\nDB Connected: %t", envOk, kiteOk, dbOk)
}

func initTickerToken() {

	envOk = loadEnv()

	if envOk {

		dbOk = db.DbInit()

		kiteOk, apiKey, accToken = kite.LoginKite()
		printStatus(envOk, dbOk, kiteOk)

		if envOk && dbOk && kiteOk {
			// Initate zerodha ticker
			kite.TickerInitialize(apiKey, accToken)
			setupCdlCrons()
			go db.StoreTickInDb()
			// start watchdog to recover from connections issues
			wdg = cron.New()
			wdg.AddFunc("@every 10s", watchdog)
			wdg.Start()

		} else {
			srv.ErrorLogger.Println("Fail to start Ticker")
		}
	}
}

func firstRunConnectionsCheck() {

	envOk = loadEnv()

	if envOk {
		dbOk = db.DbInit()
		kiteOk, apiKey, accToken = kite.LoginKite()
		printStatus(envOk, dbOk, kiteOk)
	} else {
		srv.ErrorLogger.Println("ERR: Cannot read ENV varaibles, skipping connections check!")
	}
}

func theStop() {

	kite.CloseTicker()
	cdl1min.Stop()
	cdl3min.Stop()
	cdl5min.Stop()
	wdg.Stop()
	db.CloseDBPool()
}

func setupCdlCrons() {

	cdl1min = cron.New()
	cdl1min.AddFunc("@every 1m", cdlconv.Convert1MinCandle)
	cdl1min.Start()

	cdl3min = cron.New()
	cdl3min.AddFunc("@every 3m", cdlconv.Convert3MinCandle)
	cdl3min.Start()

	cdl5min = cron.New()
	cdl5min.AddFunc("@every 5m", cdlconv.Convert5MinCandle)
	cdl5min.Start()
}

func watchdog() {
	// Initate ticker on error
	if !envOk || !dbOk || !kiteOk {
		printStatus(envOk, dbOk, kiteOk)
		theStop()
		srv.ErrorLogger.Println("\nWDG: Re-Initializing Kite", kiteOk)
		initTickerToken()
	}
}
