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

func TestStartTrader(t *testing.T) {

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	t.Parallel()

	test1(t, 1, "[case Initiate] Start two threads\n")
	test2(t, 2, "[case Initiate] daystart false, nothing should start\n")
	// test3(t, 3, "[case Resume] resume previous running trades\n")         //
	// test4(t, 4, "[case AwaitSignal] get response from api\n")             // trigger time test

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
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "$TRIGGERTIME$",
		time.Now().Local().Add(time.Second*time.Duration(10)).Format("15:04:05"), -1)

	sqlquery = strings.Replace(sqlquery, "S001-ORB", "S999-TEST", -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "PlaceOrders")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", string(appdata.ColorReset))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST ", testId, ": Trades found ", len(trades), string(appdata.ColorReset))
	}
	trades = db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", string(appdata.ColorReset))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST ", testId, ": Trades found ", len(trades), string(appdata.ColorReset))
	}
	time.Sleep(time.Second * 7)
	trades = db.ReadAllOrderBookFromDb("=", "PlaceOrders")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError))
		fmt.Print("TEST  ", testId, ": FAILED\n")
	} else {
		fmt.Print(string(appdata.ColorSuccess))
		fmt.Print("PASSED: TEST_", testId, ": Trades found ", len(trades))
	}
	// terminate trademgr
	// TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}
func test3(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_setup)
	db.DbRawExec(Test3_orderbook)
	// start trader, do not spawn new trades
	go StartTrader(false)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "TradeMonitoring")
	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n")
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
