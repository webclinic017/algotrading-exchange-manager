package db

import (
	"algo-ex-mgr/app/srv"
	"os"
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
	{"BANKNIFTY-FUT", 0, 0, "FUT", "2022-03-23", "2022-04-23", "BANKNIFTY22MARFUT"},

	{"BANKNIFTY-FUT", 35120, 0, "CE", "2022-03-23", "2022-03-30", "BANKNIFTY2232435200CE"},
	{"BANKNIFTY-FUT", 35120, 1, "CE", "2022-03-23", "2022-03-30", "BANKNIFTY2232435300CE"},
	{"BANKNIFTY-FUT", 35120, -5, "CE", "2022-03-23", "2022-03-30", "BANKNIFTY2232434700CE"},
	{"BANKNIFTY-FUT", 35120, 0, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232435200PE"},
	{"BANKNIFTY-FUT", 36020, 0, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232436100PE"},
	{"BANKNIFTY-FUT", 36020, 1, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232436200PE"},
	{"BANKNIFTY-FUT", 36020, -1, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232436000PE"},
	{"BANKNIFTY-FUT", 36020, -11, "PE", "2022-03-23", "2022-03-30", "BANKNIFTY2232435000PE"},

	{"ASHOK LEYLAND", 0, 0, "EQ", "2022-03-23", "2022-03-30", "ASHOKLEY"},
	{"RELIANCE INDUSTRIES", 0, 0, "EQ", "2022-03-23", "2022-03-30", "RELIANCE"},
}

func TestFetchInstrData(t *testing.T) {
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	DbInit()

	for _, test := range FetchInstrDataTests {

		actual, _ := FetchInstrData(test.instrument, test.strikelevel, test.opdepth, test.optype, test.startdate, test.enddate)

		if actual != test.expected {
			t.Errorf("\nderiveFuturesName() \nexpected:%q \n  actual:%q", test.expected, actual)
		}
	}
}

func TestGetInstrumentsToken(t *testing.T) {
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	DbInit()

	actual := GetInstrumentsToken()

	if len(actual) == 0 {
		t.Errorf("\nGetInstrumentsToken() \nexpected:%q \n  actual:%q", 2, len(actual))
	} else {
		println("\nGetInstrumentsToken() \nexpected:%q \n  actual:%q", 2, len(actual))
	}
}
