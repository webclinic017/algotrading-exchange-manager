package kite

import (
	"goTicker/app/data"
	"testing"
	"time"
)

type CalOrderMarginTesting struct {
	argDate, argInstr string
	argWeekSel        int
	argStrikePrice    float64
	argOptionLevel    int
	argDirection      string
	argOrderRoute     string
	expected          string
}

var CalOrderMarginTests = []CalOrderMarginTesting{}

func TestCalOrderMargin(t *testing.T) {
	var order data.TradeSignal
	var ts data.Strategies
	// var om []kiteconnect.OrderMargins

	for _, test := range DeriveOptionNameTests {

		order.Instr = test.argInstr
		order.Entry = test.argStrikePrice
		order.Dir = test.argDirection
		ts.CtrlParam.TradeSettings.OptionLevel = test.argOptionLevel
		ts.CtrlParam.TradeSettings.OptionExpiryWeek = test.argWeekSel
		ts.CtrlParam.TradeSettings.OrderRoute = test.argOrderRoute
		// expected := test.expected

		actual := CalOrderMargin(order, ts, time.Now())

		if actual[0].TradingSymbol != test.expected {
			t.Errorf("deriveFuturesName()")
		}
	}
}
