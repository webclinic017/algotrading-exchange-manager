package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"
	"testing"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func TestFinalizeOrder_LIVE(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()

	if appdata.Env["ZERODHA_LIVE_TEST"] != "TRUE" {
		t.Errorf(appdata.ErrorColor, "\n\nLive testing is disabled. Set ZERODHA_LIVE_TEST to TRUE in userSettings.env")
		return
	}

	var warCheck = "Check results MANUALLY!!!"
	type ExecuteOrderT struct {
		argStrategy    string
		argInstr       string
		argEntry       bool
		argWeekSel     int
		argStrikePrice float64
		argOptionLevel int
		argDirection   string
		argOrderRoute  string
		argMonthSel    int
		argSkipExpWk   bool
		argOrderType   string
		argVarieties   string
		argValidities  string
		argProducts    string
		orderPlaced    bool
	}

	var order appdata.OrderBook_S
	var ts appdata.UserStrategies_S

	// ** This is live testcase - update dates are per current symbols dates and levels.
	// ** Result needs to be verified manually!!!
	var testArray = []ExecuteOrderT{

		// order placed only in trading time, rejected otherwise as
		// we are using Market order (hard coded) and not limit order for options
		{"SRB-001", "BANKNIFTY-FUT", true, 0 + WeekSel, 35000 + StrikePrice, 0 + OptionLevel,
			"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

		{"SRB-001", "BANKNIFTY-FUT", false, 0 + WeekSel, 35000 + StrikePrice, 0 + OptionLevel,
			"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

		{"SRB-001", "BANKNIFTY-FUT", true, 0 + WeekSel, 35000 + StrikePrice, 0 + OptionLevel,
			"bullish", "option-sell", 0 + MonthSel, SkipExpWkTrue, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

		{"SRB-001", "BANKNIFTY-FUT", false, 0 + WeekSel, 35000 + StrikePrice, 0 + OptionLevel,
			"bullish", "option-sell", 0 + MonthSel, SkipExpWkTrue, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

		// This order shall be placed (ensure strike price is within reach)
		{"SRB-001", "ZEEL-FUT", true, 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
			"bullish", "futures", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},

		{"SRB-001", "ZEEL-FUT", false, 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
			"bullish", "futures", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},

		// This order shall be placed (ensure strike price is within reach)
		{"SRB-001", "ASHOKLEY", true, 0 + WeekSel, 140 + StrikePrice, 0 + OptionLevel,
			"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

		// This order shall not get placed, since strikeprice is 0 (outside of circuit limit)
		{"SRB-001", "ICICIPRULI", true, 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
			"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},

		// This order shall not get placed, since strikeprice is 0 (outside of circuit limit)
		{"SRB-001", "ICICIPRULI", false, 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
			"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},

		{"SRB-901", "ICICIPRULI", false, 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
			"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
			kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},
	}

	fmt.Printf(appdata.ErrorColor, warCheck)
	for _, test := range testArray {

		ts.Strategy = test.argStrategy
		ts.Parameters.Kite_Setting.Varieties = test.argVarieties
		ts.Parameters.Kite_Setting.Products = test.argProducts
		ts.Parameters.Kite_Setting.Validities = test.argValidities
		ts.Parameters.Kite_Setting.OrderType = test.argOrderType
		ts.Parameters.Futures_Setting.FuturesExpiryMonth = test.argMonthSel
		ts.Parameters.Futures_Setting.SkipExipryWeekFutures = test.argSkipExpWk
		ts.Parameters.Kite_Setting.OrderRoute = test.argOrderRoute
		ts.Parameters.Option_setting.OptionExpiryWeek = test.argWeekSel
		ts.Parameters.Option_setting.OptionLevel = test.argOptionLevel

		order.Dir = test.argDirection
		order.Instr = test.argInstr
		order.Targets.EntrPrice = test.argStrikePrice

		// expected := test.expected

		var ordModify uint64
		if test.argStrategy == "SRB-901" {
			ordModify = 1
		} else {
			ordModify = 0
		}
		orderID := finalizeOrder(order, ts, time.Now(), 1, ordModify, test.argEntry)

		if orderID == 0 && test.orderPlaced == true {
			t.Errorf(appdata.ErrorColor, "\nderiveFuturesName() No data fetched - check dates and levels are correct. This UT is live with server\n")
		} else {
			// print result for manual check
			fmt.Printf(appdata.InfoColor, order.Instr)
			fmt.Printf(appdata.InfoColor, test.argOrderRoute)
			fmt.Printf(appdata.InfoColorUint, orderID)
			fmt.Println()
		}

	}
	fmt.Println()
}

func TestGetLowestPrice(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	kite.Init()
	fmt.Print(appdata.ColorBlue, "\nThis test requires market to be open, else live quotes return 0.\n")

	type tstst struct {
		id          uint32
		instr       string
		dir         string
		orderMargin float64
	}

	var tstArry = []tstst{
		{2, "unknown", "sell", 0},
		{1, "INFY", "buy", 0},
		{2, "INFY", "sell", 0},
		{1, "RELIANCE", "buy", 0},
		{2, "BANKNIFTY22JUNFUT-FUT", "buy", 0}, // this name is time dependent, ensure you change to relevant
		{2, "BANKNIFTY22JUNFUT-FUT", "sell", 0},
	}

	for _, test := range tstArry {

		price := getLowestPrice(test.instr, test.dir)
		if price == 0 {
			t.Errorf("getLowestPrice() (Test-%v) Error fetching %v-%v", test.id, test.instr, test.dir)
		}
	}
	fmt.Println()
}

func TestDetermineOrderSize(t *testing.T) {

	type DetermineOrderSizeT struct {
		userMargin  float64
		orderMargin float64
		winningRate float64
		maxBudget   float64
		limitAmount float64
		expResult   float64
	}

	var DetermineOrderSizeTestArray = []DetermineOrderSizeT{

		// {0.0, 0.0, 0.0, 0.0, 0, 0},
		{50001.0, 40000.0, 100.0, 100.0, 0, 0},
		{50002.0, 40000.0, 100.0, 100.0, 50000, 1},
		{50002.0, 40000.0, 100.0, 100.0, 5000, 0},
		{50003.0, 40000.0, 0.0, 100.0, 50000, 1},
		{100004.0, 40000.0, 10.0, 100.0, 100004, 1},
		{100005.0, 40000.0, 100.0, 100.0, 100000, 2},
		{100006.0, 40000.0, 80.0, 100.0, 100000, 2},
		{1007.0, 4000.0, 100.0, 100.0, 1000, 0},
		{10008.0, 1000.0, 100.0, 100.0, 10000, 10},
		{10001.0, 1000.0, 100.0, 10.0, 10000, 1},
		{0.0, 1000.0, 100.0, 10.0, 10000, 0},
	}

	for _, test := range DetermineOrderSizeTestArray {

		qty := determineOrderSize(test.userMargin, test.orderMargin, test.winningRate, test.maxBudget, test.limitAmount)
		if test.expResult != qty {
			t.Errorf("determineOrderSize() usermargin:%v expected:%v  actual:%v", test.userMargin, test.expResult, qty)
		}
	}
	fmt.Println()
}

func TestPendingOrderEntr(t *testing.T) {

	fmt.Printf(appdata.ColorReset)
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	kite.Init()
	db.DbInit()

	type testSt struct {
		id            uint
		orderIdEntr   uint64
		orderIdExit   uint64
		simulation    bool
		qtyReq        float64
		qtyFilledEntr float64
		result        bool
	}

	// ** This is live testcase - update dates are per current symbols dates and levels.
	// ** Result needs to be verified manually!!!
	var tArray = []testSt{
		{1, 0, 0, true, 0, 0, true},
		{2, 0, 0, false, 1, 0, false},
		{3, 0, 0, false, 1, 1, true},
		{4, 0, 0, false, 10, 9, false},
		{5, 0, 0, false, 0, 0, true},
		{6, 220529100004459, 0, false, 0, 0, true}, // real id fetch
	}

	var order appdata.OrderBook_S
	var us appdata.UserStrategies_S

	for _, test := range tArray {

		order.Info.Order_simulation = test.simulation
		order.Info.OrderIdEntr = test.orderIdEntr
		order.Info.OrderIdExit = test.orderIdExit
		order.Info.QtyReq = test.qtyReq
		order.Info.QtyFilledEntr = test.qtyFilledEntr

		val := pendingOrderEntr(&order, us)

		if val != test.result {
			t.Errorf("pendingOrderEntr() (Test-%v) failed", test.id)
		}
	}
}

func TestPendingOrderExit(t *testing.T) {

	fmt.Printf(appdata.ColorReset)
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	kite.Init()
	db.DbInit()

	type testSt struct {
		id            uint
		orderIdEntr   uint64
		orderIdExit   uint64
		simulation    bool
		qtyFilledEntr float64
		QtyFilledExit float64
		result        bool
	}

	// ** This is live testcase - update dates are per current symbols dates and levels.
	// ** Result needs to be verified manually!!!
	var tArray = []testSt{
		{1, 0, 0, true, 0, 0, true},
		{2, 0, 0, false, 1, 0, false},
		{3, 0, 0, false, 1, 1, true},
		{4, 0, 0, false, 10, 9, false},
		{5, 0, 0, false, 0, 0, true},
		{6, 0, 220529100004459, false, 0, 0, true}, // real id fetch
	}

	var order appdata.OrderBook_S
	var us appdata.UserStrategies_S

	for _, test := range tArray {

		order.Info.Order_simulation = test.simulation
		order.Info.OrderIdEntr = test.orderIdEntr
		order.Info.OrderIdExit = test.orderIdExit
		order.Info.QtyFilledEntr = test.qtyFilledEntr
		order.Info.QtyFilledExit = test.QtyFilledExit

		val := pendingOrderExit(&order, us)

		if val != test.result {
			t.Errorf("pendingOrderEntr() (Test-%v) failed", test.id)
		}
	}
}

func TestPendingTradeEnterExit(t *testing.T) {

	fmt.Printf(appdata.ColorReset)
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	kite.Init()
	db.DbInit()

	type testSt struct {
		id            uint
		simulation    bool
		QtyFilledEntr float64
		result        bool
	}

	// ** This is live testcase - update dates are per current symbols dates and levels.
	// ** Result needs to be verified manually!!!
	var tArray = []testSt{
		{1, true, 0, true},
		{2, false, 1, false},
	}
	var order appdata.OrderBook_S
	var us appdata.UserStrategies_S

	for _, test := range tArray {

		order.Info.Order_simulation = test.simulation
		order.Info.QtyFilledEntr = test.QtyFilledEntr

		val := tradeEnter(&order, us)
		if val != test.result {
			t.Errorf("pendingOrderEntr() (Test-%v) failed", test.id)
		}
		val = tradeExit(&order, us)
		if val != test.result {
			t.Errorf("pendingOrderEntr() (Test-%v) failed", test.id)
		}

	}
}
