package main

import (
	"fmt"
	"goTicker/app/cdlconv"
	"goTicker/app/db"
	"goTicker/app/kite"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk                                     bool
	apiKey, accToken                                        string
	cdl1min, cdl3min, cdl5min, wdg, closeTicker, initTicker *cron.Cron
)

func main() {

	// timeZone, _ := time.LoadLocation("Asia/Calcutta")

	initTickerToken() // start now, when docker starts

	// everyday scheduled start
	initTicker = cron.New()
	initTicker.AddFunc("0 23 0 * * *", initTickerToken)
	initTicker.Start()

	// stop connection everyday
	closeTicker = cron.New()
	closeTicker.AddFunc("0 24 0 * * *", theStop)
	closeTicker.Start()

	// start watchdog to recover from connections issues
	wdg = cron.New()
	wdg.AddFunc("@every 10s", watchdog)
	wdg.Start()

	select {}

}

func loadEnv() bool {

	// Load .env file, if not in production
	if 0 >= len(os.Getenv("PRODUCTION")) {
		err := godotenv.Load("config/ENV_Settings.env")
		if err != nil {
			log.Fatal("ENV_Settings.env file not found, Terminating!!!")
			//return false
		}
	}
	if 0 >= len(os.Getenv("TFA_AUTH")) {
		println("TFA_AUTH not set")
		return false
	}
	if 0 >= len(os.Getenv("USER_ID")) {
		println("USER_ID not set")
		return false
	}
	if 0 >= len(os.Getenv("PASSWORD")) {
		println("PASSWORD not set")
		return false
	}
	if 0 >= len(os.Getenv("API_KEY")) {
		println("API_KEY not set")
		return false
	}
	if 0 >= len(os.Getenv("API_SECRET")) {
		println("API_SECRET not set")
		return false
	}
	if 0 >= len(os.Getenv("REQUEST_TOKEN_URL")) {
		println("REQUEST_TOKEN_URL not set")
		return false
	}
	if 0 >= len(os.Getenv("EXTERNAL_DATABASE_URL")) {
		println("EXTERNAL_DATABASE_URL not set")
		return false
	}
	os.Setenv("DATABASE_URL", os.Getenv("EXTERNAL_DATABASE_URL"))

	return true
}

func printStatus(envOk, dbOk, kiteOk bool) {
	fmt.Printf("\n--------STATUS---------")
	fmt.Printf("\nEnvironment files found: %t", envOk)
	fmt.Printf("\nKite Login Succesfull: %t", kiteOk)
	fmt.Printf("\nDB Connected: %t", dbOk)
	fmt.Printf("\n----------------------\n")
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

		} else {
			println("Fail to start Ticker")
		}
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
		fmt.Println("\nWDG: Initializing Kite", kiteOk)
		//initTickerToken()
	}

}
