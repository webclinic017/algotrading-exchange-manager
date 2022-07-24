package main

import (
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/srv"
	"os"
	"testing"
	"time"
)

func TestPostTradeOps(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	t.Parallel()

	postTradeOps()
	t.Log("Post Trade Ops Test Done")

}

func TestStartMainSession(t *testing.T) {

	startMainSession()
	time.Sleep(time.Second * 20)
	stopMainSession()
	time.Sleep(time.Second * 2)
}

func TestDbFunction(t *testing.T) {

	// appdata.ChTick = make(chan appdata.TickData, 1000)

	// _ = srv.LoadEnvVariables("app/zfiles/config/userSettings.env")
	// _ = db.DbInit()
	// go db.StoreTickInDb()
	// go kite.TestTicker()
	// println("Testing Done")

	// select {}
}

func TestTickerData(t *testing.T) {

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
