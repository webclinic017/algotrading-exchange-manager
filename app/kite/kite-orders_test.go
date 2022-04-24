package kite_test

import (
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"
	"testing"
)

func TestGetLatestQuote(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir + "/../../userSettings.env")
	// db.DbInit()
	kite.Init()

	val, n := kite.GetLatestQuote("NIFTY")

	fmt.Print(val[n].Depth.Buy[0].Price)
	fmt.Print(val[n].Depth.Buy[1].Price)
	fmt.Print(val[n].Depth.Buy[2].Price)
	fmt.Print(val[n].Depth.Buy[3].Price)
	fmt.Print(val[n].Depth.Buy[4].Price)

	// appdata.ChTick = make(chan appdata.TickData, 1000)

	// _ = srv.LoadEnvVariables("app/zfiles/config/userSettings.env")
	// _ = db.DbInit()
	// go db.StoreTickInDb()
	// go kite.TestTicker()
	// println("Testing Done")

	// select {}
}
