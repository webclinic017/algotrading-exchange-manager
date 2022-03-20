package kite

import (
	"goTicker/app/data"
	"testing"
	"time"
)

type DeriveFuturesNameTesting struct {
	argDate, argInstr string
	argMonthSel       int
	argSkipExpWk      bool
	expected          string
}

var DeriveFuturesNameTests = []DeriveFuturesNameTesting{
	{"2021-11-26", "BANKNIFTY-FUT", 0, false, "BANKNIFTY21DECFUT"},
	{"2021-11-26", "BANKNIFTY-FUT", 0, true, "BANKNIFTY21DECFUT"},
	{"2021-11-02", "BANKNIFTY-FUT", 0, true, "BANKNIFTY21NOVFUT"},
	{"2021-11-02", "BANKNIFTY-FUT", 1, false, "BANKNIFTY21DECFUT"},
	{"2021-11-02", "BANKNIFTY-FUT", 2, false, "BANKNIFTY22JANFUT"},
	{"2022-03-20", "NIFTY-FUT", 0, false, "NIFTY22MARFUT"},
	{"2022-03-24", "NIFTY-FUT", 0, false, "NIFTY22MARFUT"},
	{"2022-03-25", "NIFTY-FUT", 0, false, "NIFTY22MARFUT"},
	{"2022-03-25", "NIFTY-FUT", 0, true, "NIFTY22APRFUT"},
}

func TestDeriveFuturesName(t *testing.T) {
	var order data.TradeSignal
	var ts data.Strategies

	for _, test := range DeriveFuturesNameTests {

		dateString := test.argDate
		order.Instr = test.argInstr
		ts.CtrlParam.TradeSettings.FuturesExpiryMonth = test.argMonthSel
		ts.CtrlParam.TradeSettings.SkipExipryWeekFutures = test.argSkipExpWk
		expected := test.expected
		date, _ := time.Parse("2006-01-02", dateString)

		actual := deriveFuturesName(order, ts, date)

		if actual != expected {
			t.Errorf("deriveFuturesName() expected:%q actual:%q", expected, actual)
		}
	}
}

type DeriveOptionNameTesting struct {
	argDate, argInstr string
	argWeekSel        int
	argLevelSel       float64
	argOptnType       string
	expected          string
}

var DeriveOptionNameTests = []DeriveOptionNameTesting{
	{"2021-11-26", "BANKNIFTY-FUT", 0, 35000, "CE", "BANKNIFTY2232435000CE"},
}

func TestDeriveOptionName(t *testing.T) {
	var order data.TradeSignal
	var ts data.Strategies

	for _, test := range DeriveOptionNameTests {

		dateString := test.argDate
		order.Instr = test.argInstr
		ts.CtrlParam.TradeSettings.OptionExpiryWeek = test.argWeekSel
		order.Entry = test.argLevelSel
		expected := test.expected
		date, _ := time.Parse("2006-01-02", dateString)

		actual := deriveOptionName(order, ts, date)

		if actual != expected {
			t.Errorf("deriveFuturesName() expected:%q actual:%q", expected, actual)
		}
	}
}
