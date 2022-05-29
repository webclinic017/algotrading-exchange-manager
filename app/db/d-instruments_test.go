package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"fmt"
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

	{"ASHOKLEY", 0, 0, "EQ", "2022-03-23", "2022-03-30", "ASHOKLEY"},
	{"RELIANCE", 0, 0, "EQ", "2022-03-23", "2022-03-30", "RELIANCE"},
}

func TestFetchInstrData(t *testing.T) {
	fmt.Println(appdata.ColorInfo, "\nTestFetchInstrData()")
	fmt.Println(appdata.ColorWhite)
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	DbInit()
	DbSaveInstrCsv("user_symbols", mydir+"/../zfiles/Unittest-Support-Files/paragvb_symbols_202204282107.csv")
	DbSaveInstrCsv("instruments", mydir+"/../zfiles/Unittest-Support-Files/instruments_dbtest_data_24Mar22.csv")

	for _, test := range FetchInstrDataTests {

		actual, _ := FetchInstrData(test.instrument, test.strikelevel, test.opdepth, test.optype, test.startdate, test.enddate)

		if actual != test.expected {
			t.Error(appdata.ColorError, "\nderiveFuturesName() \nexpected:", test.expected, "\n  actual:", actual)
		} else {
			fmt.Println(appdata.ColorSuccess, "TestFetchInstrData() actual: ", actual)

		}
	}
	fmt.Println(appdata.ColorInfo)
}

func TestGetInstrumentsToken(t *testing.T) {
	fmt.Println(appdata.ColorInfo, "\nTestGetInstrumentsToken()")
	fmt.Println(appdata.ColorWhite)
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	DbInit()
	DbSaveInstrCsv("instruments", mydir+"/../zfiles/Unittest-Support-Files/instruments_dbtest_data_24Mar22.csv")

	actual := GetInstrumentsToken()

	if len(actual) == 124 {
		fmt.Println(appdata.ColorSuccess, "\nGetInstrumentsToken() expected: 124 actual: ", len(actual))
		fmt.Println(appdata.ColorInfo)
	} else {
		t.Error(appdata.ColorError, "\nGetInstrumentsToken() expected: 124 actual: ", len(actual))
		t.Error(appdata.ColorInfo)
	}
	fmt.Println(actual)

}
