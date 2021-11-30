package main

import (
	"fmt"
	"log"

	"github.com/goTicker/cdlconv"
	"github.com/goTicker/db"
	"github.com/goTicker/kite"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk            bool
	apiKey, accToken               string
	cdl1min, cdl3min, cdl5min, wdg *cron.Cron
)

func main() {

	// timeZone, _ := time.LoadLocation("Asia/Calcutta")

	// setup cron for every morning, to get fresh token
	envOk = loadEnv()
	dbOk = db.DbInit()

	initTickerToken()

	initTicker := cron.New()
	initTicker.AddFunc("0 23 0 * * *", initTickerToken)
	initTicker.Start()

	closeTicker := cron.New()
	closeTicker.AddFunc("0 24 0 * * *", theStop)
	closeTicker.Start()

	go db.StoreTickInDb()
	select {}

}

func loadEnv() bool {
	err := godotenv.Load("ENV_Settings.env", "ENV_accesstoken.env")
	if err != nil {
		log.Fatal("ENV_Settings.env / ENV_accesstoken.env file(s) not found, Terminating!!!")
		return false
	}
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

	kiteOk, apiKey, accToken = kite.LoginKite()
	printStatus(envOk, dbOk, kiteOk)
	// Do login and get access token

	if envOk && dbOk && kiteOk {
		// Initate ticker
		kite.TickerInitialize(apiKey, accToken)
		setupCdlCrons()
	} else {
		println("Fail to start Ticker")
	}
	wdg = cron.New()
	wdg.AddFunc("@every 3s", watchdog)
	wdg.Start()
}

func theStop() {

	kite.CloseTicker()
	cdl1min.Stop()
	cdl3min.Stop()
	cdl5min.Stop()
	wdg.Stop()
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
	// TODO: if kite status Nok
	// --> call initialize Kite
	fmt.Printf("\nWDG: Kite Logged in - %t", kiteOk)

}
