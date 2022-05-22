package kite

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGetLatestQuote(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	// db.DbInit()
	Init()

	val, n := GetLatestQuote("INFY")

	fmt.Print(val[n].Depth.Buy[0].Price)
	fmt.Print(val[n].Depth.Buy[1].Price)
	fmt.Print(val[n].Depth.Buy[2].Price)
	fmt.Print(val[n].Depth.Buy[3].Price)
	fmt.Print(val[n].Depth.Buy[4].Price)

	val, n = GetLatestQuote("BANKNIFTY-FUT")

	fmt.Print(appdata.ColorBlue, val[n])

	// appdata.ChTick = make(chan appdata.TickData, 1000)

	// _ = srv.LoadEnvVariables("app/zfiles/config/userSettings.env", false)
	// _ = db.DbInit()
	// go db.StoreTickInDb()
	// go kite.TestTicker()
	// println("Testing Done")

	// select {}
}

func TestFetchOrders(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	// db.DbInit()
	Init()

	tradesList := FetchOrderTrades(220425102552618)
	fmt.Println(tradesList)
	b, err := json.Marshal(tradesList)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	fmt.Println(string(b))

}
