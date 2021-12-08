package kite

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	instrument_token = iota
	exchange_token
	tradingsymbol
	name
	last_price
	expiry
	strike
	tick_size
	lot_size
	instrument_type
	segment
	exchange
)

func GetSymbols() bool {

	var (
		symbolFuturesFilter   []string
		symbolIndexFilter     []string
		symbolNseEqFilter     []string
		instrumentTokens      []string
		instrumentTokensLog   []string
		instrumentTokensError []string
	)

	e := os.Remove("config/instruments.csv")
	if e != nil {
		println("instruments.csv deleted")
	}

	fileUrl := "http://api.kite.trade/instruments"
	err := DownloadFile("config/instruments.csv", fileUrl)
	if err != nil {
		fmt.Println("Download error: instruments.csv from  " + fileUrl)
		return false
	}

	// open file
	f, err := os.Open("config/instruments.csv")
	if err != nil {
		fmt.Println("File error, cannot read instruments.csv")
		return false
	}
	// remember to close the file at the end of the program
	defer f.Close()

	csvReader := csv.NewReader(f)
	instrumentsList, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("File error, cannot read instruments.csv")
		return false
	}
	if len(instrumentsList) < 90000 {
		fmt.Println("File error, incorrect file downloaded (instruments.csv)")
		return false
	}

	// instrument_token, exchange_token,	tradingsymbol,	name
	// last_price,		 expiry,			strike,			tick_size
	// lot_size,		 instrument_type,	segment,		exchange
	//  280510214,1095743,EURINR22NOVFUT,"EURINR",0,2022-11-28,0,0.0025,1,FUT,BCD-FUT,BCD
	// fmt.Print("\n" + instrumentsList[0][instrument_token])
	// fmt.Print("\n" + instrumentsList[0][exchange_token])
	// fmt.Print("\n" + instrumentsList[0][tradingsymbol])
	// fmt.Print("\n" + instrumentsList[0][name])
	// fmt.Print("\n" + instrumentsList[0][last_price])
	// fmt.Print("\n" + instrumentsList[0][expiry])
	// fmt.Print("\n" + instrumentsList[0][strike])
	// fmt.Print("\n" + instrumentsList[0][tick_size])
	// fmt.Print("\n" + instrumentsList[0][lot_size])
	// fmt.Print("\n" + instrumentsList[0][instrument_type])
	// fmt.Print("\n" + instrumentsList[0][segment])
	// fmt.Print("\n" + instrumentsList[0][exchange] + "\n")

	dat, err := ioutil.ReadFile("config/trackSymbols.txt")
	lines := strings.Split(string(dat), "\n")
	check(err)

	symbolFuturesFilter, symbolNseEqFilter, symbolIndexFilter = sortSymbols(lines)

	iTokens, iTokensLog, iTokensError := getInstrumentTokenUniqueIdentifier(symbolFuturesFilter, instrumentsList)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	iTokens, iTokensLog, iTokensError = getInstrumentTokenNseEquity(symbolNseEqFilter, instrumentsList)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	iTokens, iTokensLog, iTokensError = getInstrumentTokenIndices(symbolIndexFilter, instrumentsList)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	saveFiles(instrumentTokens, "/log/instrumentTokens.txt")
	saveFiles(instrumentTokensLog, "log/instrumentTokensLog.txt")
	saveFiles(instrumentTokensError, "log/instrumentTokensError.txt")

	return false
}

func saveFiles(data []string, fileName string) bool {
	myfile, e := os.Create("/log/GeeksforGeeks.txt")
	if e != nil {
		log.Fatal(e)
	}
	log.Println(myfile)
	myfile.Close()

	f, err := os.Create(filepath.Base("") + fileName)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return false
	}
	// d := []string{"Welcome to the world of Go1.", "Go is a compiled language.", "It is easy to learn Go."}

	for _, v := range data {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("file written successfully")
	return true
}

func getInstrumentTokenUniqueIdentifier(symbolList []string, instrumentsList [][]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+" : "+instrumentsList[i][instrument_token])
				break
			}

		}
		if i == len(instrumentsList) {
			instrumentTokensLog = append(instrumentTokensLog, mySymbol+" : Symbol not found!")
			instrumentTokensError = append(instrumentTokensError, mySymbol+" : Symbol not found!")
		}
	}

	return instrumentTokens, instrumentTokensLog, instrumentTokensError
}

func getInstrumentTokenNseEquity(symbolList []string, instrumentsList [][]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] && "NSE" == instrumentsList[i][exchange] {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+" : "+instrumentsList[i][instrument_token])
				break
			}

		}
		if i == len(instrumentsList) {
			instrumentTokensLog = append(instrumentTokensLog, mySymbol+" : Symbol not found!")
			instrumentTokensError = append(instrumentTokensError, mySymbol+" : Symbol not found!")
		}
	}

	return instrumentTokens, instrumentTokensLog, instrumentTokensError
}

func getInstrumentTokenIndices(symbolList []string, instrumentsList [][]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] && "INDICES" == instrumentsList[i][segment] {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+" : "+instrumentsList[i][instrument_token])
				break
			}

		}
		if i == len(instrumentsList) {
			instrumentTokensLog = append(instrumentTokensLog, mySymbol+" : Symbol not found!")
			instrumentTokensError = append(instrumentTokensError, mySymbol+" : Symbol not found!")
		}
	}

	return instrumentTokens, instrumentTokensLog, instrumentTokensError
}
func sortSymbols(instrumentsList []string) ([]string, []string, []string) {
	// using for loop
	var symbolFuturesFilter []string
	var symbolIndexFilter []string
	var symbolNseEqFilter []string
	var storeIn int
	var symbolFutStr string
	const (
		noScan = iota
		futuresFilter
		nseEqFilter
		indexFilter
	)

	symbolFutStr = determineFuturesContractsName()

	for _, element := range instrumentsList {
		if strings.Contains(element, "START") {
			if strings.Contains(element, "FUTURES_FILTER") {
				storeIn = futuresFilter
				continue
			} else if strings.Contains(element, "NSEEQ_FILTER") {
				storeIn = nseEqFilter
				continue
			} else if strings.Contains(element, "INDEX_FILTER") {
				storeIn = indexFilter
				continue
			}
		} else if strings.Contains(element, "END") {
			storeIn = noScan
			continue
		}

		if storeIn == futuresFilter {
			symbolFuturesFilter = append(symbolFuturesFilter, element+symbolFutStr)
		} else if storeIn == nseEqFilter {
			symbolNseEqFilter = append(symbolNseEqFilter, element)
		} else if storeIn == indexFilter {
			symbolIndexFilter = append(symbolIndexFilter, element)
		}
	}

	return symbolFuturesFilter, symbolNseEqFilter, symbolIndexFilter
}

// func fetchInstrumentToken(symbolName string) string {

// }-

func determineFuturesContractsName() string {
	// logic -
	// 1. Jump to coming thursday
	// 2. Check if next thurday is in same month
	// 3. Use current month/year else next month/year

	var symbolFutStr string = "FAILED"
	// NIFTY21DECFUT
	dt := time.Now().Weekday()                          // todays day
	gapForThurday := math.Abs(float64(dt) - float64(4)) // 4 is thursday
	jumpToComingThurday := time.Now().AddDate(0, 0, int(gapForThurday))

	if jumpToComingThurday.Weekday().String() == "Thursday" {
		// today is Thursday

		thisMonth := time.Now().Month()
		nextWeek := time.Now().AddDate(0, 0, 7)
		monthCheck := nextWeek.Month()

		if monthCheck.String() == thisMonth.String() {
			// next thurday is in same month
			// Layouts must use the reference time Mon Jan 2 15:04:05 MST 2006 to show the pattern with which to format/parse a given time/string.
			symbolFutStr = time.Now().Format("06-Jan") + "FUT"

		} else {
			// next thurday is in next month
			symbolFutStr = nextWeek.Format("06-Jan") + "FUT"

		}
		symbolFutStr = strings.ReplaceAll(symbolFutStr, "-", "")
		symbolFutStr = strings.ToUpper(symbolFutStr)
		fmt.Println("\n\tFutures Symbol : Decoded :- ", symbolFutStr)
	}
	return symbolFutStr
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
