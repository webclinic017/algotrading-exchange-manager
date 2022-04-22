package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"

	"math"
	"testing"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

const (
	WinningRatePer = 0
	MaxBudgetPer   = 0
	LimitAmount    = 0
	WeekSel        = 0
	StrikePrice    = 0
	OptionLevel    = 0
	Direction      = 0
	OrderRoute     = 0
	MonthSel       = 0
	SkipExpWkTrue  = true
	SkipExpWkFalse = false
)

type CalOrderMarginTesting struct {
	argDate        time.Time
	argInstr       string
	argWeekSel     int
	argStrikePrice float64
	argOptionLevel int
	argDirection   string
	argOrderRoute  string
	argMonthSel    int
	argSkipExpWk   bool
	argVarieties   string
	argProducts    string
	// expected          float64
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var CalOrderMarginTests = []CalOrderMarginTesting{

	{time.Now(), "BANKNIFTY-FUT", 0 + WeekSel, 36000 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "BANKNIFTY-FUT", 0 + WeekSel, 36000 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-sell", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "BANKNIFTY-FUT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "futures", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "BANKNIFTY-FUT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bearish", "futures", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "ASHOKLEY-FUT", 0 + WeekSel, 126 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-buy", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "ASHOKLEY-FUT", 0 + WeekSel, 126 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-sell", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "ASIANPAINT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{time.Now(), "ASIANPAINT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	{time.Now(), "BHARTIARTL", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	{time.Now(), "BHARTIARTL", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	{time.Now(), "ICICIBANK-FUT", 0 + WeekSel, 762 + StrikePrice, 0 + OptionLevel,
		"bullish", "futures", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {time.Now(), "invalid", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductCNC},
}

func TestCalOrderMargin(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir + "/../../userSettings.env")
	db.DbInit()
	kite.Init()
	t.Parallel()

	var order appdata.TradeSignal
	var ts appdata.Strategies

	for _, test := range CalOrderMarginTests {

		// dateString := test.argDate
		// date, _ := time.Parse("2006-01-02", dateString)

		ts.CtrlParam.Kite_Setting.Varieties = test.argVarieties
		ts.CtrlParam.Kite_Setting.Products = test.argProducts

		ts.CtrlParam.Trade_Setting.FuturesExpiryMonth = test.argMonthSel
		ts.CtrlParam.Trade_Setting.SkipExipryWeekFutures = test.argSkipExpWk
		ts.CtrlParam.Trade_Setting.OrderRoute = test.argOrderRoute
		ts.CtrlParam.Trade_Setting.OptionExpiryWeek = test.argWeekSel
		ts.CtrlParam.Trade_Setting.OptionLevel = test.argOptionLevel

		order.Dir = test.argDirection
		order.Instr = test.argInstr
		order.Entry = test.argStrikePrice

		// expected := test.expected

		actual := getOrderMargin(order, ts, test.argDate)

		if len(actual) == 0 {
			t.Errorf(appdata.ErrorColor, "\nderiveFuturesName() No data fetched - check dates/levels/Server Auth code. This UT is live with server\n")
		} else if actual[0].Total == 0 {
			t.Errorf(appdata.ErrorColor, "\nderiveFuturesName() No margin calculated - check dates/levels/Server Auth code. This UT is live with server\n")
		} else {
			// print result for manual check

			fmt.Println()
			fmt.Printf(appdata.InfoColor, order.Instr)
			fmt.Printf(appdata.InfoColor, test.argOrderRoute)
			fmt.Printf(appdata.InfoColor, actual[0].TradingSymbol)
			fmt.Printf(appdata.InfoColorFloat, math.Round(actual[0].Total))
		}

	}
	fmt.Println()
}
