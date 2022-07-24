package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"
	"strconv"
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
	db.DbRawExec(settings_exits_deleteAll) // no exits ar defined

	subtest_StartTrader_1(t, 1, "[case Initiate] Start two threads\n")
	subtest_StartTrader_2(t, 2, "[case Initiate] daystart false, nothing should start\n")
	subtest_StartTrader_3(t, 3, "[case Resume] resume previous running trades. 1 with correct strategy set. 1 should resume\n") //
	subtest_StartTrader_4(t, 4, "[case Resume] Cannot resume - stratgey missing\n")
	subtest_StartTrader_5(t, 5, "[case Initiate] Invalid strategy\n")
	subtest_StartTrader_6(t, 6, "[case Initiate] TimeTrigged - Wait period\n")
	subtest_StartTrader_7(t, 7, "[case Initiate] TimeTrigged - Execute\n")
}
func subtest_StartTrader_1(t *testing.T, testId int, testDesc string) {

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
func subtest_StartTrader_2(t *testing.T, testId int, testDesc string) {
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
func subtest_StartTrader_3(t *testing.T, testId int, testDesc string) {
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
func subtest_StartTrader_4(t *testing.T, testId int, testDesc string) {

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
	trades := db.ReadAllOrderBookFromDb("!", "TradeCompleted")
	if len(trades) != 0 {
		t.Errorf(appdata.ColorError)
		t.Errorf("\nExpected 0 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", appdata.ColorReset)
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found ", len(trades), appdata.ColorReset, "\n")
	}
}
func subtest_StartTrader_5(t *testing.T, testId int, testDesc string) {
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
func subtest_StartTrader_6(t *testing.T, testId int, testDesc string) {

	// test if all UserStrategies are spawned
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	// setup Db entries
	TerminateTradeMgr = false
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(120)).Format("15:04:05"), -1) // after 2 min, should not execute
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-TEST-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-TEST-002", -1)
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

	StopTrader()
	time.Sleep(time.Second * 2)
}
func subtest_StartTrader_7(t *testing.T, testId int, testDesc string) {

	// test if all UserStrategies are spawned
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	// setup Db entries
	TerminateTradeMgr = false
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(3)).Format("15:04:05"), -1) // after 2 min, should not execute
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-TEST-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-TEST-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	db.DbRawExec(sqlquery)

	// start trader
	go StartTrader(true)

	time.Sleep(time.Second * 3)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "ExitTrade")
	if len(trades) != 2 {
		t.Errorf(appdata.ColorError)
		t.Errorf("\nExpected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", appdata.ColorReset)
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found ", len(trades), appdata.ColorReset, "\n")
	}

	StopTrader()
	time.Sleep(time.Second * 1)
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
	/* start trades, use active apicall, 1st trade in PlaceOrders
	modify 2nd trade time for execution, wait for timetrigger, check the second is also in PlaceOrders */

	fmt.Print((appdata.ColorReset))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()

	if appdata.Env["ZERODHA_LIVE_TEST"] != "TRUE" {
		t.Errorf(appdata.ErrorColor, "\n\nLive testing is disabled. Set ZERODHA_LIVE_TEST to TRUE in userSettings.env")
		return
	}

	fmt.Print(appdata.ColorBlue, "\nTEST_5 [case Real EQ Simulation] Simulate real equity signal and check values\n") // only at market time

	TerminateTradeMgr = false
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_EqASHOKLEY_REAL, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 3)
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d in AwaitSignal", len(trades))
		fmt.Print(appdata.ColorError, "TEST_LiveTesting : FAILED\n")
	}

	// exit trademgr - intiate sell order
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", "all-exit", -1)
	db.DbRawExec(sqlquery) // no exits ar defined
	time.Sleep(time.Second * 2)

	trades = db.ReadAllOrderBookFromDb("=", "TradeCompleted")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_LiveTesting: FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_LiveTesting: Trades found in TradeCompleted :", len(trades))
	}
}

func TestOperateSymbol_SimulationTesting(t *testing.T) {
	// Precondition
	// 1. setup user symbols
	// 2. setup user strategies

	fmt.Print((appdata.ColorReset))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	db.DbRawExec(settings_exits_deleteAll) // no exits ar defined

	testSimulation_1execute_with_termination(t, 1, "[case Initiate] Start 2 thread with 1 valid repsonce and 1 invalid resp from API to complete simulation\n")
	testSimulation_1execute_1userExit(t, 1, "[case Initiate] Start 2 thread with 1 valid repsonce and 1 invalid resp from API to complete simulation\n")
	testSimulation_1execute_with_allTermination(t, 1, "[case Initiate] Start 2 use keyword 'all-terminate, 1 trade should in Terminate'\n")
	testSimulation_1execute_with_allExit(t, 1, "[case Initiate] Start 2 use keyword 'all-exit, all trades should be TradeCompleted'\n")

}
func testSimulation_1execute_with_termination(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	TerminateTradeMgr = false
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-TEST-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 3)
	trades := db.ReadAllOrderBookFromDb("=", "ExitTrade")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in PlaceOrdersPending :", len(trades), "\nKite timeout can affect the result due to timeouts")
	}

	// terminate trademgr - trades remain in same state - no state change
	StopTrader()
	time.Sleep(time.Second * 2)

	trades = db.ReadAllOrderBookFromDb("=", "Terminate")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in Terminate state :", len(trades))
	}
}
func testSimulation_1execute_1userExit(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorReset, appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	TerminateTradeMgr = false
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-TEST-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 3)
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in AwaitSignal :", len(trades), "\nKite timeout can affect the result due to timeouts")
	}

	// terminate trademgr - trades remain in same state - no state change
	println(trades[0].Id) // trade in "AwaitSignal" state
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", strconv.FormatUint(uint64(trades[0].Id), 10), -1)
	db.DbRawExec(sqlquery) // no exits ar defined
	time.Sleep(time.Second * 3)

	trades = db.ReadAllOrderBookFromDb("=", "TradeCompleted")
	if len(trades) != 2 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in TradeCompleted state :", len(trades))
	}
}
func testSimulation_1execute_with_allTermination(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorReset, appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	TerminateTradeMgr = false
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-TEST-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 3)
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in AwaitSignal :", len(trades), "\nKite timeout can affect the result due to timeouts")
	}

	// terminate trademgr - trades remain in same state - no state change
	println(trades[0].Id) // trade in "AwaitSignal" state
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", "all-terminate", -1)
	db.DbRawExec(sqlquery) // no exits ar defined
	time.Sleep(time.Second * 3)

	trades = db.ReadAllOrderBookFromDb("=", "Terminate")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in Terminate state :", len(trades))
	}
}
func testSimulation_1execute_with_allExit(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorReset, appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	TerminateTradeMgr = false
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S999-TEST-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S999-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(2)).Format("15:04:05"), -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 2)
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 1 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in AwaitSignal :", len(trades), "\nKite timeout can affect the result due to timeouts")
	}

	println(trades[0].Id) // trade in "AwaitSignal" state
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", "all-exit", -1)
	db.DbRawExec(sqlquery) // no exits ar defined
	time.Sleep(time.Second * 3)

	trades = db.ReadAllOrderBookFromDb("=", "TradeCompleted")
	if len(trades) != 2 {
		t.Errorf("Expected 1 trades, got %d", len(trades))
		fmt.Print(appdata.ColorError, "TEST_", testId, ": FAILED\n")
	} else {
		fmt.Print(appdata.ColorSuccess, "PASSED: TEST_", testId, ": Trades found in TradeCompleted state :", len(trades))
	}
}
