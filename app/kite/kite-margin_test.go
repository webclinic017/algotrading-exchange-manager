package kite

import (
	"goTicker/app/data"
	"goTicker/app/srv"
	"testing"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
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

var CalOrderMarginTests = []CalOrderMarginTesting{
	{"2022-03-22", "INFY", 1, 73, 0, "bullish", "stock", 0, false, kiteconnect.VarietyRegular, kiteconnect.ProductMIS},
	{"2022-03-22", "BANKNIFTY-FUT", 0, 36000, 0, "bullish", "option-buy", 0, false, kiteconnect.VarietyRegular, kiteconnect.ProductNRML},
	{"2022-03-22", "BANKNIFTY-FUT", 0, 36000, 0, "bullish", "option-sell", 0, false, kiteconnect.VarietyRegular, kiteconnect.ProductNRML},
}

func TestCalOrderMargin(t *testing.T) {

	srv.Init()
	srv.LoadEnvVariables()
	SetAccessToken("y7guXmcwxu9fEb0gDHA53U9kHcshCfwB")

	var order data.TradeSignal
	var ts data.Strategies
	// var om []kiteconnect.OrderMargins

	for _, test := range CalOrderMarginTests {

		dateString := test.argDate
		date, _ := time.Parse("2006-01-02", dateString)
		order.Instr = test.argInstr
		ts.CtrlParam.TradeSettings.FuturesExpiryMonth = test.argMonthSel
		ts.CtrlParam.TradeSettings.SkipExipryWeekFutures = test.argSkipExpWk
		ts.CtrlParam.KiteSettings.Varieties = test.argVarieties
		ts.CtrlParam.KiteSettings.Products = test.argProducts

		order.Entry = test.argStrikePrice
		order.Dir = test.argDirection
		ts.CtrlParam.TradeSettings.OptionLevel = test.argOptionLevel
		ts.CtrlParam.TradeSettings.OptionExpiryWeek = test.argWeekSel
		ts.CtrlParam.TradeSettings.OrderRoute = test.argOrderRoute
		// expected := test.expected

		actual := CalOrderMargin(order, ts, date)

		if len(actual) == 0 {
			t.Errorf("deriveFuturesName() Instrument:%q Margin:%v", order.Instr, actual)
		} else if actual[0].Total == 0 {
			t.Errorf("deriveFuturesName() Instrument:%q Margin:%v", order.Instr, actual)
		} else {
			t.Logf("deriveFuturesName() Instrument:%q Margin:%v", order.Instr, actual)
		}

	}
}
