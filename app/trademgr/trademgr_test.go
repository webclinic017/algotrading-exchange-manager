package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"
	"testing"
	"time"
)

var startTrader_TblOdrbook_deleteAll = `DELETE FROM public.paragvb_order_book_test;`

var startTrader_TblStrategies_deleteAll = `DELETE FROM public.paragvb_strategies_test;`

var startTrader_TblStrategies_setup = `INSERT INTO public.paragvb_strategies_test (strategy,enabled,engine,trigger_time,trigger_days,cdl_size,instruments,controls) VALUES
('S001-ORB-001',true,'IntraDay_DNP','00:00:00','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'CONTINOUS_TT','{
"percentages": {
"target": 1,
"sl": 2,
"deepsl": 3,
"maxBudget": 50,
"winningRate":80
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
"OrderRoute": "option-buy",
"OptionLevel":-1,
"OptionExpiryWeek": 0,
"FuturesExpiryMonth": 0,    
"SkipExipryWeekFutures":true,
"LimitAmount": 30000
}
}'),
('S001-ORB-002',true,'IntraDay_DNP','22:58:00','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'TIMETRIG_TT','{
"percentages": {
"target": 1,
"sl": 2,
"deepsl": 3,
"maxBudget": 50,
"winningRate":80
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
"OrderRoute": "option-buy",
"OptionLevel":-1,
"OptionExpiryWeek": 0,
"FuturesExpiryMonth": 0,    
"SkipExipryWeekFutures":true,
"LimitAmount": 30000
}
}');`

type StartTraderT struct {
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var StartTraderTestArray = []StartTraderT{}

func TestStartTrader(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir + "/../../userSettings.env")
	db.DbInit()
	// kite.Init()
	t.Parallel()

	fmt.Print(appdata.ColorInfo, "TEST 1: Start two threads\n")
	fmt.Println(string(appdata.ColorWhite))
	// test if all strategies are spawned
	// setup Db entries
	db.DbRawExec(startTrader_TblStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)
	db.DbRawExec(startTrader_TblStrategies_setup)

	// start trader
	go StartTrader(true)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades := db.ReadAllTradeSignalFromDb("=", "AwaitSignal")
	if len(trades) != 2 {
		t.Errorf("Expected 2 trades, got %d", len(trades))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST 1: Trades found ", len(trades))
		fmt.Println(string(appdata.ColorWhite))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)

	fmt.Print(appdata.ColorInfo, "TEST 2: daystart false, nothing should start\n")
	fmt.Println(string(appdata.ColorWhite))

	db.DbRawExec(startTrader_TblStrategies_deleteAll)
	db.DbRawExec(startTrader_TblOdrbook_deleteAll)
	db.DbRawExec(startTrader_TblStrategies_setup)
	// start trader
	go StartTrader(false)

	time.Sleep(time.Second * 5)
	// check if trades are logged in order_book
	trades = db.ReadAllTradeSignalFromDb("=", "AwaitSignal")
	if len(trades) != 0 {
		t.Errorf("Expected 0 trades, got %d", len(trades))
	} else {
		fmt.Print(string(appdata.ColorSuccess), "PASSED: TEST 2: Trades found ", len(trades))
		fmt.Println(string(appdata.ColorWhite))
	}

	// terminate trademgr
	TerminateTradeMgr = true
	time.Sleep(time.Second * 3)

}
