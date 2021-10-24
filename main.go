package main

import (
	"fmt"
	"log"
	"time"

	db "github.com/goTicker/db"
	kite "github.com/goTicker/kite"
	"github.com/joho/godotenv"
)

var (
	envOk, dbOk, kiteOk bool
)

func main() {

	envOk = loadEnv()
	dbOk = db.DbInit()
	kiteOk, apiKey, accToken := kite.LoginKite()
	printStatus(envOk, dbOk, kiteOk)
	// Do login and get access token

	if envOk && dbOk && kiteOk {
		// Initate ticker
		kite.TickerInitialize(apiKey, accToken)

		time.Sleep(5 * time.Second)
		kite.CloseTicker()
		time.Sleep(5 * time.Second)
	} else {
		println("Fail to start Ticker")
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

func printStatus(envOk, dbOk, kiteOk bool) {
	fmt.Printf("\n--------STATUS---------")
	fmt.Printf("\nEnvironment files found: %t", envOk)
	fmt.Printf("\nKite Login Succesfull: %t", kiteOk)
	fmt.Printf("\nDB Connected: %t", dbOk)
	fmt.Printf("\n----------------------\n")
}
