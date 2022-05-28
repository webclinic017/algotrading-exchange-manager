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

func TestStartTrader(t *testing.T) {

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()

	test5(t, 5, "[case Initiate] Invalid strategy\n")
	test1(t, 1, "[case Initiate] Start two threads\n")
	test2(t, 2, "[case Initiate] daystart false, nothing should start\n")
	test3(t, 3, "[case Resume] resume previous running trades. 1 with correct strategy set. 1 should resume\n") //
	test4(t, 4, "[case Resume] Cannot resume - stratgey missing\n")
}

func test1(t *testing.T, testId int, testDesc string) {

	// test if all UserStrategies are spawned
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	// setup Db entries
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	db.DbRawExec(sqlquery)

	// start trader
	go StartTrader(true)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 2 {
		t.Errorf(appdata.ColorError)
		t.Errorf("\nExpected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", appdata.ColorReset)
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found ", len(trades), appdata.ColorReset, "\n")
	}
}

func test2(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	db.DbRawExec(sqlquery)

	// start trader -day false
	go StartTrader(false)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", appdata.ColorReset)
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", testId, ": Trades found ", len(trades), appdata.ColorReset)
	}
}

func test3(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continous trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)

	db.DbRawExec(sqlquery)

	time.Sleep(time.Second * 1)
	trades := db.ReadAllOrderBookFromDb("=", "ExitOrdersPending")
	if len(trades) != 0 {
		t.Errorf("\nCheck 3.1 - Expected 0 trades, got %d", len(trades))
	}
	// setup old order
	db.DbRawExec(Test3_orderbook)
	time.Sleep(time.Second * 1)
	trades = db.ReadAllOrderBookFromDb("=", "ExitOrdersPending")
	if len(trades) != 2 {
		t.Errorf("\nCheck 3.2 - Expected 2 trades, got %d", len(trades))
	}

	// start trader, do not spawn new trades
	go StartTrader(false)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades = db.ReadAllOrderBookFromDb("!=", "ExitOrdersPending")
	println(len(trades))
	if len(trades) != 0 { // len can be 0 when since no trade data from kite is stored, db parsing is resulting 0 trades.
		t.Errorf("\nCheck 3.3 - Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST_", testId, ": FAILED\n", "Check if API Server is running")

	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", testId, ": Trades found ", len(trades))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}

func test4(t *testing.T, testId int, testDesc string) {

	// test if all UserStrategies are spawned
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	// setup Db entries
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "", -1) // no days for trigger trading
	db.DbRawExec(sqlquery)

	// start trader
	go StartTrader(true)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades := db.ReadAllActiveOrderBookFromDb()
	if len(trades) != 0 {
		t.Errorf(appdata.ColorError)
		t.Errorf("\nExpected 0 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", appdata.ColorReset)
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found ", len(trades), appdata.ColorReset, "\n")
	}
}

func test5(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continous trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S000-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S000-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)

	db.DbRawExec(sqlquery)

	time.Sleep(time.Second * 1)
	trades := db.ReadAllOrderBookFromDb("=", "ExitOrdersPending")
	if len(trades) != 0 {
		t.Errorf("\nCheck 3.1 - Expected 0 trades, got %d", len(trades))
	}
	// setup old order
	db.DbRawExec(Test3_orderbook)
	time.Sleep(time.Second * 1)
	trades = db.ReadAllOrderBookFromDb("=", "ExitOrdersPending")
	if len(trades) != 2 {
		t.Errorf("\nCheck 3.2 - Expected 2 trades, got %d", len(trades))
	}

	// start trader, do not spawn new trades
	go StartTrader(false)

	time.Sleep(time.Second * 2)
	// check if trades are logged in order_book
	trades = db.ReadAllOrderBookFromDb("!=", "ExitOrdersPending")
	println(len(trades))
	if len(trades) != 0 { // len can be 0 when since no trade data from kite is stored, db parsing is resulting 0 trades.
		t.Errorf("\nCheck 3.3 - Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST_", testId, ": FAILED\n", "Check if API Server is running")

	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", testId, ": Trades found ", len(trades))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}

func TestCheckTriggerDays(t *testing.T) {

	type tst struct {
		id     int
		days   string
		result bool
		today  string
	}

	// these unit testcase are sensitive to date in "instruments" table,
	// load the instruments_dbtest_data_24Mar22.csv data before running the test cases
	var testData = []tst{
		{1, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SATURDAY"},
		{2, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SUNDAY"},
		{3, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "MONDAY"},
		{4, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "TUESDAY"},
		{5, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "WEDNESDAY"},
		{6, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", false, "THURSDAY"},
		{7, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "FRIDAY"},
		{8, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SATURDaY"},
		{9, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SundaY"},
		{10, "SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SATURDAY"},
		{11, "", false, "SATURDAY"},
		{12, "MONDAY", false, ""},
		{13, "SAturDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SATURDaY"},
		{14, "saturday, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, FRIDAY", true, "SATURDaY"},
	}

	for _, test := range testData {

		if test.result != checkTriggerDays(test.today, test.days) {
			t.Error(appdata.ColorError, " ID: ", test.id, " expected:", test.result)
		}
	}
	fmt.Println(appdata.ColorInfo)

}

func TestStartTrader_LiveTesting(t *testing.T) {

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()

	if appdata.Env["ZERODHA_LIVE_TEST"] != "TRUE" {
		t.Errorf(appdata.ErrorColor, "\n\nLive testing is disabled. Set ZERODHA_LIVE_TEST to TRUE in userSettings.env")
		return
	}
	// test4(t, 4, "[case AwaitSignal] get response from api\n")
	LiveTesting1(t, 5, "[case Real EQ Simulation] Simulate real equity signal and check values\n") // only at market time
	// test5(t, 5, "[case UserExitReq] Trade shall exit position\n")

}

/* start trades, use active apicall, 1st trade in PlaceOrders
modify 2nd trade time for execution, wait for timetrigger, check the second is also in PlaceOrders */
func LiveTesting1(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_EqASHOKLEY_REAL, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 1 {
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
