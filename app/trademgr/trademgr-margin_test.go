package trademgr

import (
	"fmt"
	"goTicker/app/data"
	"goTicker/app/db"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"math"
	"testing"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

const (
	InfoColorFloat = "\033[1;34m%.0f\033[0m\t"
	InfoColorUint  = "\033[1;34m%d\033[0m\t"
	InfoColor      = "\033[1;34m%20s\033[0m\t"
	ErrorColor     = "\033[1;31m%s\033[0m"
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
	argDate, argInstr string
	argWeekSel        int
	argStrikePrice    float64
	argOptionLevel    int
	argDirection      string
	argOrderRoute     string
	argMonthSel       int
	argSkipExpWk      bool
	argVarieties      string
	argProducts       string
	// expected          float64
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var CalOrderMarginTests = []CalOrderMarginTesting{

	// {"2022-03-22", "INFY", 1, 73, 0, "bullish", "equity", 0, false, kiteconnect.VarietyRegular, kiteconnect.ProductMIS},
	{"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 36000 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 36000 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-sell", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "futures", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bearish", "futures", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "ASHOKLEY-FUT", 0 + WeekSel, 126 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-buy", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "ASHOKLEY-FUT", 0 + WeekSel, 126 + StrikePrice, 0 + OptionLevel,
		"bullish", "option-sell", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "ASHOK LEYLAND", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	{"2022-04-02", "ASHOK LEYLAND", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	{"2022-04-02", "RELIANCE INDUSTRIES", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	{"2022-04-02", "RELIANCE INDUSTRIES", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
		kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	// {"2022-04-02", "invalid", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductCNC},
}

func TestCalOrderMargin(t *testing.T) {

	srv.Init()
	srv.LoadEnvVariables()
	db.DbInit()
	kite.SetAccessToken("qtejqw2D7IWWD0pWqMLolHzsCkWwugQf")
	t.Parallel()

	var order data.TradeSignal
	var ts data.Strategies

	for _, test := range CalOrderMarginTests {

		dateString := test.argDate
		date, _ := time.Parse("2006-01-02", dateString)

		ts.CtrlParam.KiteSettings.Varieties = test.argVarieties
		ts.CtrlParam.KiteSettings.Products = test.argProducts

		ts.CtrlParam.TradeSettings.FuturesExpiryMonth = test.argMonthSel
		ts.CtrlParam.TradeSettings.SkipExipryWeekFutures = test.argSkipExpWk
		ts.CtrlParam.TradeSettings.OrderRoute = test.argOrderRoute
		ts.CtrlParam.TradeSettings.OptionExpiryWeek = test.argWeekSel
		ts.CtrlParam.TradeSettings.OptionLevel = test.argOptionLevel

		order.Dir = test.argDirection
		order.Instr = test.argInstr
		order.Entry = test.argStrikePrice

		// expected := test.expected

		actual := CalOrderMargin(order, ts, date)

		if len(actual) == 0 {
			t.Errorf(ErrorColor, "\nderiveFuturesName() No data fetched - check dates and levels are correct. This UT is live with server\n")
		} else if actual[0].Total == 0 {
			t.Errorf(ErrorColor, "\nderiveFuturesName() No margin calculated - check dates and levels are correct. This UT is live with server\n")
		} else {
			// print result for manual check

			fmt.Println()
			fmt.Printf(InfoColor, order.Instr)
			fmt.Printf(InfoColor, test.argOrderRoute)
			fmt.Printf(InfoColor, actual[0].TradingSymbol)
			fmt.Printf(InfoColorFloat, math.Round(actual[0].Total))
		}

	}
	fmt.Println()
}
