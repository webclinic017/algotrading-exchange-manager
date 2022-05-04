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

type StartTraderT struct {
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var StartTraderTestArray = []StartTraderT{}

func TestStartTrader1(t *testing.T) {

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	t.Parallel()

	test1(t, 1, "[case Initiate] Start two threads\n")
	test2(t, 2, "[case Initiate] daystart false, nothing should start\n")
	test3(t, 3, "[case Resume] resume previous running trades. 1 with correct strategy set. 1 should resume\n") //

}

func TestStartTrader2(t *testing.T) {

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	t.Parallel()

	// test4(t, 4, "[case AwaitSignal] get response from api\n")
	test5(t, 5, "[case Real EQ Simulation] Simulate real equity signal and check values\n")
	// test5(t, 5, "[case UserExitReq] Trade shall exit position\n")

}

/* start trades,
use active apicall
1st trade in PlaceOrders
modify 2nd trade time for execution, wait for timetrigger, check the second is also in PlaceOrders
*/
func test5(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_EqRelianceREAL, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 200)
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

/* start trades,
use active apicall
1st trade in PlaceOrders
modify 2nd trade time for execution, wait for timetrigger, check the second is also in PlaceOrders
*/
func test4(t *testing.T, testId int, testDesc string) {
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
func test3(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	db.DbRawExec(sqlquery)
	db.DbRawExec(Test3_orderbook)
	// start trader, do not spawn new trades
	go StartTrader(false)

	time.Sleep(time.Second * 4)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 1 { // 1 should remain in AwaitSignal, other should be processed
		t.Errorf("Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST_", testId, ": FAILED\n", "Check if API Server is running")

	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", testId, ": Trades found ", len(trades))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}
func test1(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)
	// test if all UserStrategies are spawned
	// setup Db entries
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	db.DbRawExec(sqlquery)

	// start trader
	go StartTrader(true)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n")
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", testId, ": Trades found ", len(trades))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 1)
}
func test2(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	db.DbRawExec(sqlquery)

	// start trader
	go StartTrader(false)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n")
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", testId, ": Trades found ", len(trades))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 1)
}
