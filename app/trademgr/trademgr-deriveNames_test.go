package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/srv"
	"os"
	"testing"
	"time"
)

type DeriveInstrumentsNameTesting struct {
	argDate, argInstr string
	argWeekSel        int
	argMonthSel       int
	argStrikePrice    float64
	argOptionLevel    int
	argSkipExpWk      bool
	argDirection      string
	argOrderRoute     string
	expected          string
}

// these unit testcase are sensitive to data in "instruments" table,
// load the instruments_dbtest_data_24Mar22.csv data before running the test cases
var DeriveInstrumentsNameTests = []DeriveInstrumentsNameTesting{
	// option testing - individual securities

	{"2022-04-03", "RELIANCE-FUT", 0, 0, 2400, 0, false, "bullish", "option-buy", "RELIANCE22APR2400CE"},
	{"2022-04-12", "ASHOKLEY-FUT", 0, 0, 151.3, 0, false, "bullish", "option-buy", "ASHOKLEY22APR152.5CE"},
	{"2022-03-22", "ASHOKLEY-FUT", 1, 0, 152, 0, false, "bullish", "option-buy", "ASHOKLEY22MAR152.5CE"},
	{"2022-03-22", "ASHOKLEY-FUT", 1, 0, 152, 0, false, "bearish", "option-buy", "ASHOKLEY22MAR152.5PE"},
	// option testing - indices
	{"2022-03-25", "BANKNIFTY-FUT", 0, 0, 35123, 0, false, "bullish", "option-buy", "BANKNIFTY22MAR35200CE"},
	{"2022-03-25", "BANKNIFTY-FUT", 0, 0, 32023, 0, false, "bullish", "option-buy", "BANKNIFTY22MAR32100CE"},
	{"2022-03-25", "BANKNIFTY-FUT", 1, 0, 35123, 0, false, "bullish", "option-buy", "BANKNIFTY2240735200CE"},
	{"2022-03-25", "BANKNIFTY-FUT", 2, 0, 35123, 0, false, "bullish", "option-buy", "BANKNIFTY2241335200CE"},
	{"2022-03-25", "BANKNIFTY-FUT", 0, 0, 30023, 0, false, "bullish", "option-buy", "BANKNIFTY22MAR30100CE"},
	{"2022-03-25", "BANKNIFTY-FUT", 0, 0, 35123, 0, false, "bullish", "option-sell", "BANKNIFTY22MAR35200PE"},
	{"2022-04-04", "BANKNIFTY-FUT", 0, 0, 32089, 0, false, "bullish", "option-sell", "BANKNIFTY2240732100PE"},
	{"2022-04-04", "BANKNIFTY-FUT", 0, 0, 32089, 3, false, "bullish", "option-sell", "BANKNIFTY2240732400PE"},
	{"2022-04-04", "BANKNIFTY-FUT", 0, 0, 32089, -3, false, "bullish", "option-sell", "BANKNIFTY2240731800PE"},
	{"2022-04-15", "NIFTY-FUT", 0, 0, 16123, 0, false, "bullish", "option-sell", "NIFTY2242116150PE"},
	{"2022-04-15", "NIFTY-FUT", 0, 0, 16123, 0, false, "BEARISh", "option-sell", "NIFTY2242116150CE"},
	{"2022-04-15", "NIFTY-FUT", 0, 0, 16123, 0, false, "bullish", "option-buy", "NIFTY2242116150CE"},
	{"2022-04-15", "NIFTY-FUT", 1, 0, 16123, 0, false, "bullish", "option-buy", "NIFTY22APR16150CE"},
	{"2022-04-15", "NIFTY-FUT", 3, 0, 16123, 0, false, "buLLish", "option-buy", "NIFTY2251216150CE"},
	{"2022-04-15", "NIFTY-FUT", 1, 0, 16000, 0, false, "bullIsh", "option-buy", "NIFTY22APR16000CE"},
	{"2022-04-15", "NIFTY-FUT", 1, 0, 16036, 0, false, "bullish", "option-buy", "NIFTY22APR16050CE"},
	// futues testing - indices
	{"2022-04-15", "NIFTY-FUT", 0, 0, 16036, 0, false, "bullish", "futures", "NIFTY22APRFUT"},
	{"2022-04-15", "NIFTY-FUT", 0, 1, 0, 0, false, "bullish", "futures", "NIFTY22MAYFUT"},
	{"2022-04-15", "NIFTY-FUT", 0, 1, 0, 0, false, "bearish", "futures", "NIFTY22MAYFUT"},
	// futures testing - individual securities
	{"2022-04-15", "ASHOKLEY-FUT", 0, 1, 0, 0, false, "bullish", "futures", "ASHOKLEY22MAYFUT"},
	{"2022-04-15", "ASHOKLEY-FUT", 0, 0, 0, 0, false, "bullish", "futures", "ASHOKLEY22APRFUT"},
	{"2022-04-15", "RELIANCE-FUT", 0, 1, 0, 0, false, "bullish", "futures", "RELIANCE22MAYFUT"},
	{"2022-03-25", "RELIANCE-FUT", 0, 0, 0, 0, false, "bullish", "futures", "RELIANCE22MARFUT"},
	// Equite testing - individual securities
	{"2022-04-15", "ASHOKLEY", 0, 0, 0, 0, false, "", "equity", "ASHOKLEY"},
	{"2020-04-15", "ASHOKLEY", 0, 0, 0, 0, false, "", "equity", "ASHOKLEY"},
	{"2022-04-15", "RELIANCE", 0, 0, 0, 0, false, "", "equity", "RELIANCE"},
}

func TestDeriveInstrumentsName(t *testing.T) {
	t.Parallel()
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	db.DbInit()
	db.DbSaveInstrCsv("instruments", mydir+"/../zfiles/Unittest-Support-Files/instruments_dbtest_data_24Mar22.csv")

	var order appdata.OrderBook_S
	var ts appdata.UserStrategies_S

	for _, test := range DeriveInstrumentsNameTests {

		dateString := test.argDate
		date, _ := time.Parse("2006-01-02", dateString)
		order.Instr = test.argInstr
		order.Targets.EntrPrice = test.argStrikePrice
		order.Dir = test.argDirection
		ts.Parameters.Option_setting.OptionLevel = test.argOptionLevel
		ts.Parameters.Futures_Setting.FuturesExpiryMonth = test.argMonthSel
		ts.Parameters.Futures_Setting.SkipExipryWeekFutures = test.argSkipExpWk
		ts.Parameters.Option_setting.OptionExpiryWeek = test.argWeekSel
		ts.Parameters.Kite_Setting.OrderRoute = test.argOrderRoute
		expected := test.expected

		actual, _ := deriveInstrumentsName(order, ts, date)

		if actual != expected {
			t.Errorf("deriveName() \nexpected:%q \n  actual:%q", expected, actual)
		}
	}
}
