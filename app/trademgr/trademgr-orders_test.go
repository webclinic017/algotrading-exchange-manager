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

var warCheck = "Check results MANUALLY!!!"

type ExecuteOrderT struct {
	argStrategy    string
	argInstr       string
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

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var ExecuteOrderTestArray = []ExecuteOrderT{

	// order placed only in trading time, rejected otherwise
	{"SRB-001", "BANKNIFTY-FUT", 0 + WeekSel, 40000 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue, kiteconnect.OrderTypeLimit,
		kiteconnect.VarietyRegular, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},

	// This order shall be placed (ensure strike price is within reach)
	{"SRB-001", "ICICIPRULI", 0 + WeekSel, 540 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
		kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

	// This order shall not get placed, since strikeprice is 0 (outside of circuit limit)
	{"SRB-001", "ICICIPRULI", 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
		kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},
}

func TestPlaceOrder(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	t.Parallel()

	var order appdata.OrderBook_S
	var ts appdata.UserStrategies_S

	fmt.Printf(appdata.ErrorColor, warCheck)
	for _, test := range ExecuteOrderTestArray {

		ts.Strategy = test.argStrategy
		ts.Parameters.Kite_Setting.Varieties = test.argVarieties
		ts.Parameters.Kite_Setting.Products = test.argProducts
		ts.Parameters.Kite_Setting.Validities = test.argValidities
		ts.Parameters.Kite_Setting.OrderType = test.argOrderType
		ts.Parameters.Futures_Setting.FuturesExpiryMonth = test.argMonthSel
		ts.Parameters.Futures_Setting.SkipExipryWeekFutures = test.argSkipExpWk
		ts.Parameters.Option_setting.OrderRoute = test.argOrderRoute
		ts.Parameters.Option_setting.OptionExpiryWeek = test.argWeekSel
		ts.Parameters.Option_setting.OptionLevel = test.argOptionLevel

		order.Dir = test.argDirection
		order.Instr = test.argInstr
		order.Targets.Entry = test.argStrikePrice

		// expected := test.expected

		orderID := executeOrder(order, ts, time.Now(), 1)

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

func TestDetermineOrderSize(t *testing.T) {

	for _, test := range DetermineOrderSizeTestArray {

		qty := determineOrderSize(test.userMargin, test.orderMargin, test.winningRate, test.maxBudget, test.limitAmount)
		if test.expResult != qty {
			t.Errorf("determineOrderSize() usermargin:%v expected:%v  actual:%v", test.userMargin, test.expResult, qty)
		}

	}
	fmt.Println()
}

type EnterTradeT struct {
	argStrategy    string
	argInstr       string
	MaxBudget      float64
	WinningRate    float64
	LimitAmount    float64
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

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var EnterTradeTestArray = []EnterTradeT{

	// // order placed only in trading time, rejected otherwise
	// {"SRB-001", "BANKNIFTY-FUT", 0 + WeekSel, 40000 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue, kiteconnect.OrderTypeLimit,
	// 	kiteconnect.VarietyRegular, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},

	// This order shall be placed (ensure strike price is within reach)
	{"SRB-001", "ICICIPRULI", 100 + MaxBudgetPer, 100 + WinningRatePer, 1000 + LimitAmount, 0 + WeekSel, 540 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
		kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

	// {"SRB-001", "ICICIPRULI", 50 + MaxBudgetPer, 100 + WinningRatePer, 1000 + LimitAmount, 0 + WeekSel, 540 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
	// 	kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, true},

	// // This order shall not get placed, since strikeprice is 0 (outside of circuit limit)
	// {"SRB-001", "ICICIPRULI", 50 + MaxBudgetPer, 100 + WinningRatePer, 1000 + LimitAmount, 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
	// 	kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS, false},
}

func TestEnterTrade(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	kite.Init()
	t.Parallel()

	var order appdata.OrderBook_S
	var ts appdata.UserStrategies_S

	fmt.Printf(appdata.ErrorColor, warCheck)
	for _, test := range EnterTradeTestArray {

		ts.Strategy = test.argStrategy
		ts.Parameters.Kite_Setting.Varieties = test.argVarieties
		ts.Parameters.Kite_Setting.Products = test.argProducts
		ts.Parameters.Kite_Setting.Validities = test.argValidities
		ts.Parameters.Kite_Setting.OrderType = test.argOrderType
		ts.Parameters.Futures_Setting.FuturesExpiryMonth = test.argMonthSel
		ts.Parameters.Futures_Setting.SkipExipryWeekFutures = test.argSkipExpWk
		ts.Parameters.Option_setting.OrderRoute = test.argOrderRoute
		ts.Parameters.Option_setting.OptionExpiryWeek = test.argWeekSel
		ts.Parameters.Option_setting.OptionLevel = test.argOptionLevel
		ts.Parameters.Controls.MaxBudget = test.MaxBudget
		ts.Parameters.Controls.WinningRatio = test.WinningRate
		ts.Parameters.Controls.LimitAmount = test.LimitAmount

		order.Dir = test.argDirection
		order.Instr = test.argInstr
		order.Targets.Entry = test.argStrikePrice

		// expected := test.expected

		// tradeEnter(order, ts)

		// if orderID == 0 && test.orderPlaced == true {
		// 	t.Errorf(ErrorColor, "\nderiveFuturesName() No data fetched - check dates and levels are correct. This UT is live with server\n")
		// } else {
		// print result for manual check
		fmt.Printf(appdata.InfoColor, order.Instr)
		fmt.Printf(appdata.InfoColor, test.argOrderRoute)
		// fmt.Printf(InfoColorUint, orderID)
		fmt.Println()
		// }

	}
	fmt.Println()
}
