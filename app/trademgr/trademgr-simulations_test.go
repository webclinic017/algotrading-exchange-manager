package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStartTrader_SimulationTesting(t *testing.T) {
	// Precondition
	// 1. setup user symbols
	// 2. setup user strategies

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	t.Parallel()

	testSimulation(t, 1, "[case Initiate] Start two threads\n")

}

func testSimulation(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 2 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", string(appdata.ColorReset))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST ", testId, ": Trades found in AwaitSignal", len(trades), string(appdata.ColorReset))
	}
	time.Sleep(time.Second * 9)
	trades = db.ReadAllOrderBookFromDb("=", "PlaceOrdersPending")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in PlaceOrdersPending", len(trades))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)

	trades = db.ReadAllOrderBookFromDb("=", "Terminate")
	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in Terminate", len(trades))
	}
}
