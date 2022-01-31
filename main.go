package main

import (
	"goTicker/app/db"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"goTicker/app/trademgr"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk           bool
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

	startSession(false) // Check if conections are okay

	// start watchdog to recover from connections issues
	wdg = cron.New()
	wdg.AddFunc("@every 30s", checkConnection)
	wdg.Start()

	// everyday scheduled start At 09:00:00 Mon-Fri
	initTicker = cron.New()
	initTicker.AddFunc("0 0 9 * * 1-5", startKite)
	initTicker.Start()

	// everyday scheduled stop At 16:00:00 Mon-Fri
	closeTicker = cron.New()
	closeTicker.AddFunc("0 0 16 * * 1-5", stopKite)
	closeTicker.Start()

	select {}

}

func loadEnv() bool {

	// Load .env file, if not in production

	println("PRODUCTION - ", os.Getenv("PRODUCTION"))
	if os.Getenv("PRODUCTION") != "true" {
		srv.WarningLogger.Println("DEVELOPMENT ENV")
		srv.InfoLogger.Println("Ensure ENV variables are set in ENV_settings.env")
		srv.FileCopyIfMissing("app/templates/ENV_Settings.env", "app/config/ENV_Settings.env")
		_ = godotenv.Load("app/config/ENV_Settings.env")
	} else {
		srv.InfoLogger.Println("PRODUCTION ENV- Ensure ENV variables are set")
	}

	srv.InfoLogger.Println("user ID", os.Getenv("USER_ID"))

	if 0 >= len(os.Getenv("LIVE_TRADING_MODE")) {
		srv.ErrorLogger.Println("LIVE_TRADING_MODE not set")
		return false
	}

	if 0 >= len(os.Getenv("USER_ID")) {
		srv.ErrorLogger.Println("USER_ID not set")
		return false
	}
	if 0 >= len(os.Getenv("TFA_AUTH")) {
		srv.ErrorLogger.Println("TFA_AUTH not set")
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
	if 0 >= len(os.Getenv("TIMESCALEDB_ADDRESS")) {
		srv.ErrorLogger.Println("TIMESCALEDB_ADDRESS not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_USERNAME")) {
		srv.ErrorLogger.Println("TIMESCALEDB_USERNAME not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_PASSWORD")) {
		srv.ErrorLogger.Println("TIMESCALEDB_PASSWORD not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_PORT")) {
		srv.ErrorLogger.Println("TIMESCALEDB_PORT not set")
		return false
	}
	return true
}

func startSession(check bool) {

	srv.InfoLogger.Println("\n~~~~~~~~~~~~~~~~~~~~~~~~ Let's begin -", time.Now().Format("Monday, Jan-02 3:4 PM"), "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	envOk = loadEnv()

	if envOk {

		dbOk = db.DbInit()
		if dbOk {
			kite.Tokens, kite.InsNamesMap, symbolFutStr, symbolMcxFutStr = kite.GetSymbols()
			db.StoreSymbolsInDb(symbolFutStr, symbolMcxFutStr)
		}

		kiteOk, apiKey, accToken = kite.LoginKite()

		if envOk && dbOk && kiteOk && check {
			// Initate zerodha ticker
			kite.TickerInitialize(apiKey, accToken)
			go db.StoreTickInDb()
			go trademgr.Trader()

			// start watchdog to recover from connections issues
		}
	}
}
func stopKite() {
	kite.CloseTicker()
}

func startKite() {
	kite.CloseTicker()
	srv.InfoLogger.Println("\nInitializing Kite", kiteOk)
	startSession(true)
	srv.InfoLogger.Printf("\n\n\t--------------STATUS---------------\n\t| Environment variables set: %t |\n\t| Kite Login Succesfull: %t     |\n\t| DB Connected: %t              |\n\t-----------------------------------\n\n", envOk, kiteOk, dbOk)
}

func checkConnection() {

	if !kite.KiteConnectionStatus {
		now := time.Now()
		if (now.Hour() >= 9) && (now.Hour() < 16) &&
			(now.Weekday() > 0) && (now.Weekday() < 6) {
			startKite()
		}
	}
}

func testDbFunction() {

	kite.ChTick = make(chan kite.TickData, 1000)

	_ = loadEnv()
	_ = db.DbInit()
	go db.StoreTickInDb()
	go kite.TestTicker()
	println("Testing Done")

	select {}
}

func testTickerData() {

	startKite()
	time.Sleep(time.Second * 20)
	stopKite()
	time.Sleep(time.Second * 20)
	startKite()
	time.Sleep(time.Second * 20)
	stopKite()
	println("Testing Done")

	select {}
}
