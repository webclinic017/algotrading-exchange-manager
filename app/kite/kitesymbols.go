package kite

import (
	"encoding/csv"
	"fmt"
	"goTicker/app/srv"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const ( // format in which data is received form zerodha server
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

func GetSymbols() ([]uint32, map[string]string, string, string) {

	var (
		symbolNseFuturesFilter []string
		symbolMcxFuturesFilter []string
		symbolIndexFilter      []string
		symbolNseEqFilter      []string
		instrumentTokens       []string
		instrumentTokensLog    []string
		instrumentTokensError  []string
		instrumentUint32       []uint32
	)

	srv.InfoLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Fetch Symbols",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	insMap := make(map[string]string)

	e := os.Remove("app/log/instruments.csv")
	if e != nil {
		srv.InfoLogger.Println("instruments.csv deleted")
	}

	fileUrl := "http://api.kite.trade/instruments"
	err := DownloadFile("app/log/instruments.csv", fileUrl)
	if err != nil {
		srv.ErrorLogger.Println("Download error: instruments.csv from  " + fileUrl)
		return instrumentUint32, insMap, symbolFutStr, symbolMcxFutStr
	}

	// open file
	f, err := os.Open("app/log/instruments.csv")
	if err != nil {
		srv.ErrorLogger.Println("File error, cannot read instruments.csv")
		return instrumentUint32, insMap, symbolFutStr, symbolMcxFutStr
	}
	// remember to close the file at the end of the program
	defer f.Close()

	csvReader := csv.NewReader(f)
	instrumentsList, err := csvReader.ReadAll()
	if err != nil {
		srv.ErrorLogger.Println("File error, cannot read instruments.csv")
		return instrumentUint32, insMap, symbolFutStr, symbolMcxFutStr
	}
	if len(instrumentsList) < 90000 {
		srv.ErrorLogger.Println("File error, incorrect file downloaded (instruments.csv)")
		return instrumentUint32, insMap, symbolFutStr, symbolMcxFutStr
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

	dat, err := ioutil.ReadFile("app/config/trackSymbols.txt")
	lines := strings.Split(string(dat), "\n")
	check(err)

	symbolFutStr = determineFuturesContractsName()
	symbolMcxFutStr = determineMcxFuturesContractsName()

	symbolNseFuturesFilter, symbolMcxFuturesFilter, symbolNseEqFilter, symbolIndexFilter = sortSymbols(lines, symbolFutStr, symbolMcxFutStr)

	iTokens, iTokensLog, iTokensError := getInstrumentTokenNseFutures(symbolNseFuturesFilter, instrumentsList, insMap)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	iTokens, iTokensLog, iTokensError = getInstrumentTokenMcxFutures(symbolMcxFuturesFilter, instrumentsList, insMap)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	iTokens, iTokensLog, iTokensError = getInstrumentTokenNseEquity(symbolNseEqFilter, instrumentsList, insMap)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	iTokens, iTokensLog, iTokensError = getInstrumentTokenIndices(symbolIndexFilter, instrumentsList, insMap)
	instrumentTokens = append(instrumentTokens, iTokens...)
	instrumentTokensLog = append(instrumentTokensLog, iTokensLog...)
	instrumentTokensError = append(instrumentTokensError, iTokensError...)

	saveFiles(instrumentTokens, "instrumentTokensList.log")
	saveFiles(instrumentTokensLog, "instrumentTokensFound.log")
	saveFiles(instrumentTokensError, "instrumentTokensParseError.log")

	srv.ErrorLogger.Println(instrumentTokensError)

	return convertStringArrayToUint32Array(instrumentTokens), insMap, symbolFutStr, symbolMcxFutStr
}

func convertStringArrayToUint32Array(symbolList []string) []uint32 {

	var symbolListUint32 []uint32

	for _, mySymbol := range symbolList {
		val, err := strconv.Atoi(mySymbol)
		if err != nil {
			srv.ErrorLogger.Println("\nError converting string to uint32")
		}
		symbolListUint32 = append(symbolListUint32, uint32(val))
	}
	return symbolListUint32
}

func saveFiles(data []string, fileName string) bool {
	// logic
	// 1. delete if file exists
	// 2. create file
	// 3. write data to file
	// 4. close file

	e := os.Remove("app/log/" + fileName)
	if e != nil {
		srv.InfoLogger.Println(fileName + " deleted")
	}

	f, err := os.Create("app/log/" + fileName)

	if err != nil {
		srv.ErrorLogger.Println(err)
		f.Close()
		return false
	}

	fmt.Fprintln(f, "File generated at : "+time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range data {
		fmt.Fprintln(f, v)
		if err != nil {
			srv.ErrorLogger.Println(err)
			return false
		}
	}
	err = f.Close()
	if err != nil {
		srv.ErrorLogger.Println(err)
		return false
	}
	return true
}

func getInstrumentTokenNseFutures(symbolList []string, instrumentsList [][]string, insMap map[string]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+","+instrumentsList[i][instrument_token])
				insFlatName := strings.ReplaceAll(mySymbol, symbolFutStr, "")
				insFlatName = strings.ReplaceAll(insFlatName, symbolMcxFutStr, "")
				insMap[instrumentsList[i][instrument_token]] = insFlatName + "-FUT"
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

func getInstrumentTokenMcxFutures(symbolList []string, instrumentsList [][]string, insMap map[string]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+","+instrumentsList[i][instrument_token])
				insFlatName := strings.ReplaceAll(mySymbol, symbolFutStr, "")
				insFlatName = strings.ReplaceAll(insFlatName, symbolMcxFutStr, "")
				insMap[instrumentsList[i][instrument_token]] = insFlatName + "-MCXFUT"
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

func getInstrumentTokenNseEquity(symbolList []string, instrumentsList [][]string, insMap map[string]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] && instrumentsList[i][exchange] == "NSE" {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+","+instrumentsList[i][instrument_token])
				insMap[instrumentsList[i][instrument_token]] = mySymbol
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

func getInstrumentTokenIndices(symbolList []string, instrumentsList [][]string, insMap map[string]string) ([]string, []string, []string) {

	var instrumentTokens []string
	var instrumentTokensLog []string
	var instrumentTokensError []string
	var i int

	for _, mySymbol := range symbolList {

		for i = 0; i < len(instrumentsList); i++ {
			if mySymbol == instrumentsList[i][tradingsymbol] && instrumentsList[i][segment] == "INDICES" {
				instrumentTokens = append(instrumentTokens, instrumentsList[i][instrument_token])
				instrumentTokensLog = append(instrumentTokensLog, mySymbol+","+instrumentsList[i][instrument_token])
				insMap[instrumentsList[i][instrument_token]] = mySymbol + "-IDX"
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
func sortSymbols(instrumentsList []string, symbolFutStr string, symbolMcxFutStr string) ([]string, []string, []string, []string) {
	// using for loop
	var symbolNseFuturesFilter []string
	var symbolMcxFuturesFilter []string
	var symbolIndexFilter []string
	var symbolNseEqFilter []string
	var storeIn int
	// var symbolFutStr string
	// var symbolMcxFutStr string
	const (
		noScan = iota
		nseFuturesFilter
		mcxFuturesFilter
		nseEqFilter
		indexFilter
	)

	for _, element := range instrumentsList {
		if strings.Contains(element, "START") {
			if strings.Contains(element, "NSE_FUTURES") {
				storeIn = nseFuturesFilter
				continue
			} else if strings.Contains(element, "MCX_FUTURES") {
				storeIn = mcxFuturesFilter
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

		if storeIn == nseFuturesFilter {
			symbolNseFuturesFilter = append(symbolNseFuturesFilter, element+symbolFutStr)
		} else if storeIn == mcxFuturesFilter {
			symbolMcxFuturesFilter = append(symbolMcxFuturesFilter, element+symbolMcxFutStr)
		} else if storeIn == nseEqFilter {
			symbolNseEqFilter = append(symbolNseEqFilter, element)
		} else if storeIn == indexFilter {
			symbolIndexFilter = append(symbolIndexFilter, element)
		}
	}

	return symbolNseFuturesFilter, symbolMcxFuturesFilter, symbolNseEqFilter, symbolIndexFilter
}

func determineMcxFuturesContractsName() string {
	// logic -
	// 1. Jump to coming thursday
	// 2. Check if next thurday is in same month
	// 3. Use current month/year else next month/year

	var symbolFutStr string = "FAILED"
	var jumpToNextContract time.Time

	mnt := time.Now().Month() // current month
	if mnt == time.February || mnt == time.April || mnt == time.June || mnt == time.August || mnt == time.October || mnt == time.December {
		jumpToNextContract = time.Now().AddDate(0, 2, 0) // jump to next contract, two months ahead
	} else {
		jumpToNextContract = time.Now().AddDate(0, 1, 0)
	}

	symbolFutStr = jumpToNextContract.Format("06-Jan") + "FUT"
	symbolFutStr = strings.ReplaceAll(symbolFutStr, "-", "")
	symbolFutStr = strings.ToUpper(symbolFutStr)
	srv.InfoLogger.Println("MCX Futures Symbol : Decoded :- ", symbolFutStr)

	return symbolFutStr
}

func determineFuturesContractsName() string {
	// logic -
	// 1. Jump to coming thursday
	// 2. Check if next thurday is in same month
	// 3. Use current month/year else next month/year

	var symbolFutStr string = "FAILED"
	// NIFTY21DECFUT
	dt := time.Now().Weekday() // todays day
	gapForThurday := 0

	if dt == time.Monday {
		gapForThurday = 3
	} else if dt == time.Tuesday {
		gapForThurday = 2
	} else if dt == time.Wednesday {
		gapForThurday = 1
	} else if dt == time.Friday {
		gapForThurday = 6
	} else if dt == time.Saturday {
		gapForThurday = 5
	} else if dt == time.Sunday {
		gapForThurday = 4
	} else {
		gapForThurday = 0 // its thursday, do nothing
	}

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
		srv.InfoLogger.Println("Futures Symbol : Decoded :- ", symbolFutStr)
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
