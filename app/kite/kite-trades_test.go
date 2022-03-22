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
		date, _ := time.Parse("2006-01-02", dateString)
		order.Instr = test.argInstr
		ts.CtrlParam.TradeSettings.FuturesExpiryMonth = test.argMonthSel
		ts.CtrlParam.TradeSettings.SkipExipryWeekFutures = test.argSkipExpWk
		expected := test.expected

		actual := deriveFuturesName(order, ts, date)

		if actual != expected {
			t.Errorf("deriveFuturesName() expected:%q actual:%q", expected, actual)
		}
	}
}

type DeriveOptionNameTesting struct {
	argDate, argInstr string
	argWeekSel        int
	argStrikePrice    float64
	argOptionLevel    int
	argDirection      string
	argOrderRoute     string
	expected          string
}

var DeriveOptionNameTests = []DeriveOptionNameTesting{
	{"2022-03-22", "ASHOKLEY-FUT", 1, 73, 0, "bullish", "option-buy", "ASHOKLEY22MAR72.5CE"},
	{"2022-03-22", "TVSMOTOR-FUT", 1, 484, 0, "bullish", "option-buy", "TVSMOTOR22MAR480CE"},
	{"2022-03-22", "TVSMOTOR-FUT", 1, 500, 0, "bullish", "option-buy", "TVSMOTOR22MAR500CE"},
	{"2022-03-22", "TVSMOTOR-FUT", 1, 484, 5, "bullish", "option-buy", "TVSMOTOR22MAR530CE"},
	{"2022-03-22", "TVSMOTOR-FUT", 1, 484, -5, "bullish", "option-buy", "TVSMOTOR22MAR430CE"},
	{"2022-03-22", "TVSMOTOR-FUT", 1, 484, -5, "BEARISH", "option-buy", "TVSMOTOR22MAR530PE"},
	{"2022-03-22", "TVSMOTOR-FUT", 1, 484, -5, "bullish", "option-sell", "TVSMOTOR22MAR530PE"},
	{"2022-03-22", "ASHOKLEY-FUT", 1, 25, 0, "bullish", "option-buy", "ASHOKLEY22MAR25CE"},
	{"2022-03-22", "ASHOKLEY-FUT", 1, 25.8, 0, "bullish", "option-buy", "ASHOKLEY22MAR25CE"},
	{"2022-03-15", "BANKNIFTY-FUT", 0, 35123, 0, "bullish", "option-buy", "BANKNIFTY2231735100CE"},
	{"2022-03-15", "BANKNIFTY-FUT", 0, 32023, 0, "bullish", "option-buy", "BANKNIFTY2231732000CE"},
	{"2022-03-15", "BANKNIFTY-FUT", 1, 35123, 0, "bullish", "option-buy", "BANKNIFTY2232435100CE"},
	{"2022-03-15", "BANKNIFTY-FUT", 2, 35123, 0, "bullish", "option-buy", "BANKNIFTY22MAR35100CE"},
	{"2022-03-15", "BANKNIFTY-FUT", 0, 35123, 0, "bullish", "option-buy", "BANKNIFTY2231735100CE"},
	{"2022-03-15", "BANKNIFTY-FUT", 0, 35123, 0, "bullish", "option-buy", "BANKNIFTY2231735100CE"},
	{"2021-12-22", "BANKNIFTY-FUT", 0, 35123, 0, "bullish", "option-buy", "BANKNIFTY21D2335100CE"},
	{"2021-12-22", "BANKNIFTY-FUT", 0, 35123, 0, "bullish", "option-sell", "BANKNIFTY21D2335100PE"},
	{"2021-12-22", "NIFTY-FUT", 0, 35123, 0, "bullish", "option-sell", "NIFTY21D2335100PE"},
	{"2021-12-22", "NIFTY-FUT", 0, 35123, 0, "BEARISh", "option-sell", "NIFTY21D2335100CE"},
	{"2021-12-22", "NIFTY-FUT", 0, 35123, 0, "bullish", "option-buy", "NIFTY21D2335100CE"},
	{"2022-03-22", "NIFTY-FUT", 1, 35123, 0, "bullish", "option-buy", "NIFTY22MAR35100CE"},
	{"2022-03-22", "NIFTY-FUT", 3, 35123, 0, "bullish", "option-buy", "NIFTY2241435100CE"},
	{"2022-03-22", "NIFTY-FUT", 1, 35000, 0, "bullish", "option-buy", "NIFTY22MAR35000CE"},
	{"2022-03-22", "NIFTY-FUT", 1, 2536, 0, "bullish", "option-buy", "NIFTY22MAR2500CE"},
}

func TestDeriveOptionName(t *testing.T) {
	var order data.TradeSignal
	var ts data.Strategies

	for _, test := range DeriveOptionNameTests {

		dateString := test.argDate
		date, _ := time.Parse("2006-01-02", dateString)
		order.Instr = test.argInstr
		order.Entry = test.argStrikePrice
		order.Dir = test.argDirection
		ts.CtrlParam.TradeSettings.OptionLevel = test.argOptionLevel
		ts.CtrlParam.TradeSettings.OptionExpiryWeek = test.argWeekSel
		ts.CtrlParam.TradeSettings.OrderRoute = test.argOrderRoute
		expected := test.expected

		actual := deriveOptionName(order, ts, date)

		if actual != expected {
			t.Errorf("deriveFuturesName() expected:%q actual:%q", expected, actual)
		}
	}
}
