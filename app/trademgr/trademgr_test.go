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

var startTrader_TblOdrbook_deleteAll = `DELETE FROM public.paragvb_order_book_test;`

var startTrader_TblUserStrategies_deleteAll = `DELETE FROM public.paragvb_strategies_test;`

var startTrader_TblUserStrategies_setup = `INSERT INTO public.paragvb_strategies_test (strategy,enabled,engine,trigger_time,trigger_days,cdl_size,instruments,controls) VALUES
('S001-ORB-001',true,'IntraDay_DNP','00:00:00','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'ASHOKLEY','{
"percentages": {
"target": 11,
"sl": 21,
"deepsl": 31,
"maxBudget": 51,
"winningRate":81
},
"target_controls": {
"trail_target_en": false,
"position_reversal_en": true,
"delayed_stoploss_min": "2018-09-22T23:23:23Z",
"stall_detect_period_min": "2018-09-22T22:22:22Z"
},
"kite_setting": {
"products": "MIS",
"varieties": "regular",
"OrderType": "MARKET",
"Validities": "IOC",
"PositionType": "day"
},
"trade_setting": {
"OrderRoute": "equity",
"OptionLevel":-1,
"OptionExpiryWeek": 0,
"FuturesExpiryMonth": 0,    
"SkipExipryWeekFutures":true,
"LimitAmount": 30000
}
}'),
('S001-ORB-002',true,'IntraDay_DNP','$TRIGGERTIME$','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'ASHOKLEY','{
"percentages": {
"target": 12,
"sl": 22,
"deepsl": 32,
"maxBudget": 52,
"winningRate":82
},
"target_controls": {
"trail_target_en": false,
"position_reversal_en": true,
"delayed_stoploss_min": "2018-09-22T23:23:23Z",
"stall_detect_period_min": "2018-09-22T22:22:22Z"
},
"kite_setting": {
"products": "MIS",
"varieties": "regular",
"OrderType": "MARKET",
"Validities": "IOC",
"PositionType": "day"
},
"trade_setting": {
"OrderRoute": "equity",
"OptionLevel":-1,
"OptionExpiryWeek": 0,
"FuturesExpiryMonth": 0,    
"SkipExipryWeekFutures":true,
"LimitAmount": 30000
}
}');`

var test3_orderbook = `INSERT INTO public.paragvb_order_book_test ("date",instr,strategy,status,instr_id,dir,entry,target,stoploss,order_id,order_trade_entry,order_trade_exit,order_simulation,exit_reason,post_analysis) VALUES
('2022-04-20','CONTINOUS_test1','S001-ORB-001','TradeMonitoring',0,'',0.0,0.0,0.0,0,'{}','{}','{}','','{}'),
('2022-04-20','TIMETRIG_test2','S001-ORB-002','TradeMonitoring',0,'',0.0,0.0,0.0,0,'{}','{}','{}','','{}');`

type StartTraderT struct {
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var StartTraderTestArray = []StartTraderT{}

func TestStartTrader(t *testing.T) {

	fmt.Print((appdata.ColorWhite))
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir + "/../../userSettings.env")
	db.DbInit()
	kite.Init()
	t.Parallel()

	test4(t, 4)

	// test1(t, 1)
	// test2(t, 2)
	// test3(t, 3)

}

/* start trades,
use active apicall
1st trade in PlaceOrders
modify 2nd trade time for execution, wait for timetrigger, check the second is also in PlaceOrders
*/
func test4(t *testing.T, testId int) {
	fmt.Print(appdata.ColorInfo, "\nTEST  : [case AwaitSignal] get response from api", string(appdata.ColorDimmed))

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)

	// add 10 seconds to timetriggered trade
	sqlquery := strings.Replace(startTrader_TblUserStrategies_setup, "$TRIGGERTIME$",
		time.Now().Local().Add(time.Second*time.Duration(10)).Format("15:04:05"), -1)

	sqlquery = strings.Replace(sqlquery, "S001-ORB", "S999-TEST", -1)

	db.DbRawExec(sqlquery)

	go StartTrader(true)

	time.Sleep(time.Second * 10)
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
		fmt.Print("TEST  ", testId, ": FAILED\n", string(appdata.ColorDimmed))
	} else {
		fmt.Print(string(appdata.ColorSuccess))
		fmt.Print("PASSED: TEST ", testId, ": Trades found ", len(trades), string(appdata.ColorReset))
	}
	// terminate trademgr
	// TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}
func test3(t *testing.T, testId int) {
	fmt.Print(appdata.ColorInfo, "\nTEST  ", testId, ": [case Resume] resume previous running trades")
	fmt.Println(string(appdata.ColorWhite))

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_setup)
	db.DbRawExec(test3_orderbook)
	// start trader, do not spawn new trades
	go StartTrader(false)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "TradeMonitoring")
	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", string(appdata.ColorWhite))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST  ", testId, ": Trades found ", len(trades), string(appdata.ColorReset))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}

func test1(t *testing.T, testId int) {

	fmt.Print(appdata.ColorInfo, "\nTEST  ", testId, ": [case Initiate] Start two threads\n", string(appdata.ColorWhite))
	// test if all UserStrategies are spawned
	// setup Db entries
	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_setup)

	// start trader
	go StartTrader(true)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", string(appdata.ColorWhite))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST  ", testId, ": Trades found ", len(trades), string(appdata.ColorWhite))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}

func test2(t *testing.T, testId int) {
	fmt.Print(appdata.ColorInfo, "\nTEST  ", testId, ": [case Initiate] daystart false, nothing should start\n")
	fmt.Println(string(appdata.ColorWhite))

	db.DbRawExec(startTrader_TblUserStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)
	db.DbRawExec(startTrader_TblUserStrategies_setup)
	// start trader
	go StartTrader(false)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllOrderBookFromDb("=", "AwaitSignal")
	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d", len(trades))
		fmt.Print((appdata.ColorError), "TEST  ", testId, ": FAILED\n", string(appdata.ColorWhite))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST  ", testId, ": Trades found ", len(trades), string(appdata.ColorReset))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)
}
