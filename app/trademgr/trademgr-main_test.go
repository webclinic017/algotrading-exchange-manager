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

func _exit() {
	// terminate trademgr
	StopTrader()
	time.Sleep(time.Second * 3)
	fmt.Print((appdata.ColorInfo), "\n")
}
func _checkOrderBook(t *testing.T, tname string, s uint64, condition string, value string, expVal int) {

	time.Sleep(time.Second * time.Duration(s))
	trades := db.ReadAllOrderBookFromDb(condition, value)

	if len(trades) != expVal {
		t.Errorf("\nCheck "+tname+" - Expected %d trades, got %d", expVal, len(trades))
		fmt.Print((appdata.ColorError), "TEST_", tname, ": FAILED\n", "Check if API Server is running")
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST_", tname, ": Trades found in ", value, "-", len(trades), "\n")
	}
}
func _resetOrderBook(sql string) string {
	// setup Db entries
	TerminateTradeMgr = false
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// make continour trigerred trades
	sqlquery := strings.Replace(sql, "%TRIGGERTIME", "00:00:00", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_1", "S990-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S990-CONT-002", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_1", "TT_TEST1", -1)
	sqlquery = strings.Replace(sqlquery, "%SYMBOL_NAME_2", "TT_TEST2", -1)
	sqlquery = strings.Replace(sqlquery, "%TRIGGER_DAYS", "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday", -1)

	return sqlquery
}

// #################################################################################################### CheckTriggerDays
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

// #################################################################################################### StartTrader - LIVE
func TestStartTrader_LiveTesting_AfterMarketOnly(t *testing.T) {
	/* start trades, use active apicall, 1st trade in PlaceOrders
	modify 2nd trade time for execution, wait for timetrigger, check the second is also in PlaceOrders */

	fmt.Print((appdata.ColorReset))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	TerminateTradeMgr = false

	if appdata.Env["ZERODHA_LIVE_TEST"] != "TRUE" {
		t.Errorf(appdata.ErrorColor, "\n\nLive testing is disabled. Set ZERODHA_LIVE_TEST to TRUE in userSettings.env")
		return
	}
	fmt.Print(appdata.ColorBlue, "\nTEST [case Real EQ Trade] Place real equity order (after market) and check values\n") // only at market time
	db.DbRawExec(settings_exits_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	db.DbRawExec(_resetOrderBook(startTrader_TblUserStrategies_EqASHOKLEY_REAL))

	go StartTrader(true)
	_checkOrderBook(t, "TestStartTrader_LiveTesting_AfterMarketOnly", 20, "=", "PlaceOrdersPending", 1)

	// exit trademgr - intiate sell order
	sqlquery := strings.Replace(settings_exits_setVal, "%EXIT_ID", "all-exit", -1)
	db.DbRawExec(sqlquery) // no exits ar defined

	_checkOrderBook(t, "TestStartTrader_LiveTesting_ExitTrade", 20, "=", "TradeCompleted", 1)
	_exit()

}

// #################################################################################################### StartTrader
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
	subtest_StartTrader_4(t, 4, "[case Resume] Cannot resume - strategy day not enabled\n")
	subtest_StartTrader_5(t, 5, "[case Invalid Strategy] Invalid strategy. New and Resume\n")
	subtest_StartTrader_6(t, 6, "[case Wait for Trigger] TimeTrigged - Wait period\n")
	subtest_StartTrader_7(t, 7, "[case Execute on Trigger] TimeTrigged - Completee Execution\n")
}
func subtest_StartTrader_1(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)
	db.DbRawExec(_resetOrderBook(startTrader_TblUserStrategies_setup))

	go StartTrader(true)
	_checkOrderBook(t, "1.1", 2, "=", "AwaitSignal", 2)
}
func subtest_StartTrader_2(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)
	db.DbRawExec(_resetOrderBook(startTrader_TblUserStrategies_setup))

	go StartTrader(false) // start trader -day false
	_checkOrderBook(t, "2.1", 2, "=", "AwaitSignal", 0)
}
func subtest_StartTrader_3(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)
	db.DbRawExec(_resetOrderBook(startTrader_TblUserStrategies_setup))
	_checkOrderBook(t, "3.1", 1, "=", "ExitOrdersPending", 0)

	// setup old order
	db.DbRawExec(Test3_orderbook)
	_checkOrderBook(t, "3.2", 2, "=", "ExitOrdersPending", 2) // check if the trades are in orderbook

	go StartTrader(false) // now start trader, do not spawn new trades
	_checkOrderBook(t, "3.3", 4, "=", "TradeCompleted", 2)

	_exit()
}
func subtest_StartTrader_4(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGER_DAYS", "", -1) // no days for trigger trading
	sqlquery = _resetOrderBook(sqlquery)
	db.DbRawExec(sqlquery)

	go StartTrader(true)
	_checkOrderBook(t, "4.1", 3, "!=", "TradeCompleted", 0)

	_exit()
}
func subtest_StartTrader_5(t *testing.T, testId int, testDesc string) {
	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%STRATEGY_NAME_1", "S000-CONT-001", -1)
	sqlquery = strings.Replace(sqlquery, "%STRATEGY_NAME_2", "S000-CONT-002", -1)
	sqlquery = _resetOrderBook(sqlquery)
	db.DbRawExec(sqlquery)

	db.DbRawExec(sqlquery)
	_checkOrderBook(t, "5.1", 1, "=", "ExitOrdersPending", 0)

	db.DbRawExec(Test3_orderbook) // setup old order
	_checkOrderBook(t, "5.2", 1, "=", "ExitOrdersPending", 2)

	go StartTrader(false)
	_checkOrderBook(t, "5.3", 2, "!=", "ExitOrdersPending", 0)

	_exit()
}
func subtest_StartTrader_6(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Second*time.Duration(120)).Format("15:04:05"), -1) // after 2 min, should not execute
	db.DbRawExec(_resetOrderBook(sqlquery))

	// start trader
	TerminateTradeMgr = false
	go StartTrader(true)
	_checkOrderBook(t, "6.1", 2, "=", "AwaitSignal", 2)

	_exit()
}
func subtest_StartTrader_7(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)
	TerminateTradeMgr = false

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Minute*time.Duration(1)).Format("15:04:05"), -1) // execute after 1 min
	db.DbRawExec(_resetOrderBook(sqlquery))

	go StartTrader(true)
	_checkOrderBook(t, "7.1", 75, "=", "TradeCompleted", 2)

	_exit()
}

// #################################################################################################### StartTrader - SIMULATION
func TestOperateSymbol_TerminateTrades(t *testing.T) {
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

	subtest_TerminateTrades_1(t, 1, "[case - Terminate all using API StopTrades() ]\n")
	subtest_TerminateTrades_2(t, 1, "[case - Exit thread based on ID ]\n")
	subtest_TerminateTrades_3(t, 1, "[case - all-exit  from db settings ]\n")
	subtest_TerminateTrades_4(t, 1, "[case - all-terminate  from db settings ]\n")
}
func subtest_TerminateTrades_1(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Minute*time.Duration(1)).Format("15:04:05"), -1) // execute after 1 min
	db.DbRawExec(_resetOrderBook(sqlquery))

	go StartTrader(true)
	_checkOrderBook(t, "Terminate 1.1", 5, "=", "AwaitSignal", 2)

	// terminate trademgr - trades remain in same state - no state change
	StopTrader()
	_checkOrderBook(t, "Terminate 1.2", 10, "=", "Terminate", 2)

	_exit()
}
func subtest_TerminateTrades_2(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorReset, appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Minute*time.Duration(1)).Format("15:04:05"), -1) // execute after 1 min
	db.DbRawExec(_resetOrderBook(sqlquery))

	go StartTrader(true)
	_checkOrderBook(t, "Terminate 2.1", 5, "=", "AwaitSignal", 2)

	// read id and set in db for exiting that trade
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	println(trades[0].Id)
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", strconv.FormatUint(uint64(trades[0].Id), 10), -1)
	db.DbRawExec(sqlquery)

	_checkOrderBook(t, "Terminate 2.1", 20, "=", "TradeCompleted", 1)
	_exit()

}
func subtest_TerminateTrades_3(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorReset, appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Minute*time.Duration(1)).Format("15:04:05"), -1) // execute after 1 min
	db.DbRawExec(_resetOrderBook(sqlquery))

	go StartTrader(true)
	_checkOrderBook(t, "Terminate 3.1", 5, "=", "AwaitSignal", 2)

	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	println(trades[0].Id)
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", "all-exit", -1)
	db.DbRawExec(sqlquery)

	_checkOrderBook(t, "Terminate 3.1", 20, "=", "TradeCompleted", 2)
	_exit()

}
func subtest_TerminateTrades_4(t *testing.T, testId int, testDesc string) {

	fmt.Print(appdata.ColorReset, appdata.ColorBlue, "\nTEST_", testId, ": ", testDesc)

	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "%TRIGGERTIME",
		time.Now().Local().Add(time.Minute*time.Duration(1)).Format("15:04:05"), -1) // execute after 1 min
	db.DbRawExec(_resetOrderBook(sqlquery))

	go StartTrader(true)
	_checkOrderBook(t, "Terminate 4.1", 5, "=", "AwaitSignal", 2)

	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	println(trades[0].Id)
	sqlquery = strings.Replace(settings_exits_setVal, "%EXIT_ID", "all-terminate", -1)
	db.DbRawExec(sqlquery)

	_checkOrderBook(t, "Terminate 4.1", 20, "=", "Terminate", 2)
	_exit()

}
