package db

import (
	"goTicker/app/srv"
	"testing"
)

type FetchInstrDataTesting struct {
	instrument  string
	strikelevel uint64
	opdepth     int
	optype      string
	startdate   string
	enddate     string
	expected    string
}

// these unit testcase are sensitive to date in "instruments" table,
// load the instruments_dbtest_data_24Mar22.csv data before running the test cases
var FetchInstrDataTests = []FetchInstrDataTesting{
	{"BANKNIFTY-FUT", 35120, 0, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232435200PE"},
	{"BANKNIFTY-FUT", 36020, 0, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232436100PE"},
	{"BANKNIFTY-FUT", 36020, 1, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232436200PE"},
	{"BANKNIFTY-FUT", 36020, -1, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232436000PE"},
	{"BANKNIFTY-FUT", 36020, -11, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232435000PE"},
}

func TestFetchInstrData(t *testing.T) {
	srv.Init()
	srv.LoadEnvVariables()
	DbInit()

	for _, test := range FetchInstrDataTests {

		actual, _ := FetchInstrData(test.instrument, test.strikelevel, test.opdepth, test.optype, test.startdate, test.enddate)

		if actual != test.expected {
			t.Errorf("\nderiveFuturesName() \nexpected:%q \n  actual:%q", test.expected, actual)
		}
	}
}
