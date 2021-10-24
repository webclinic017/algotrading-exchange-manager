package main

import (
	"fmt"
	"log"
	"time"

	"github.com/goTicker/dbticks"
	"github.com/goTicker/kite"
	"github.com/joho/godotenv"
)

var (
	envOk, dbOk, kiteOk bool
)

func main() {

	envOk = loadEnv()
	dbOk = dbticks.DbInit()
	kiteOk, apiKey, accToken := kite.LoginKite()

	// Do login and get access token

	if kiteOk && envOk && dbOk {
		// Initate ticker
		kite.TickerInitialize(apiKey, accToken)

		time.Sleep(5 * time.Second)
		kite.CloseTicker()
		time.Sleep(5 * time.Second)
	} else {
		println("Fail to start Ticker")
		fmt.Printf("\nEnvironment files found: %t", envOk)
		fmt.Printf("\nKite Login Succesfull: %t", kiteOk)
		fmt.Printf("\nDB Connected: %t", dbOk)
	}
}

func loadEnv() bool {
	err := godotenv.Load("ENV_Settings.env", "ENV_accesstoken.env")
	if err != nil {
		log.Fatal("ENV_Settings.env / ENV_accesstoken.env file(s) not found, Terminating!!!")
		return false
	}
	return true
}
