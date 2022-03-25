package trademgr

import (
	"fmt"
	"goTicker/app/data"
	"goTicker/app/db"

	"strings"
	"time"
)

// The format is BANKNIFTY<YY><M><DD>strike<PE/CE>
// The month format is 1 for JAN, 2 for FEB, 3, 4, 5, 6, 7, 8, 9, O(capital o) for October, N for November, D for December.
// var symbolFutStr string = "FAILED"
// BANKNIFTY2232435000CE - 24th Mar 2022
// BANKNIFTY22MAR31000CE - 31st Mar 2022
// Last week of Month - will be monthly expiry
func deriveOptionName(order data.TradeSignal, ts data.Strategies, selDate time.Time) string {

	var (
		instrumentType string
		strStartDate   string
		strEndDate     string
		enddate        time.Time
	)

	// ---------------------------------------------------------------------- COMPUTE CE/PE
	if ts.CtrlParam.TradeSettings.OrderRoute == "option-buy" {
		enddate = selDate.AddDate(0, 0, 7+(7*ts.CtrlParam.TradeSettings.OptionExpiryWeek))
		if strings.ToLower(order.Dir) == "bullish" {
			instrumentType = "CE"
		} else {
			instrumentType = "PE"
		}
	} else if ts.CtrlParam.TradeSettings.OrderRoute == "option-sell" {
		enddate = selDate.AddDate(0, 0, 7+(7*ts.CtrlParam.TradeSettings.OptionExpiryWeek))
		if strings.ToLower(order.Dir) == "bullish" {
			instrumentType = "PE"
		} else {
			instrumentType = "CE"
		}
	} else if ts.CtrlParam.TradeSettings.OrderRoute == "futures" {
		enddate = selDate.AddDate(0, 0, 30+(30*ts.CtrlParam.TradeSettings.FuturesExpiryMonth))
		instrumentType = "FUT"
	}

	selDate = selDate.AddDate(0, 0, (7 * ts.CtrlParam.TradeSettings.OptionExpiryWeek))
	strStartDate = fmt.Sprintf("%d-%02d-%02d", selDate.Year(), selDate.Month(), selDate.Day())

	// ---------------------------------------------------------------------- Special case for expirt
	// For individual securities expiry is monthly
	if (strings.ToLower(order.Instr) != "nifty-fut") ||
		(strings.ToLower(order.Instr) != "banknifty-fut") ||
		(strings.ToLower(order.Instr) != "finnifty-fut") ||
		(strings.ToLower(order.Instr) != "midcpnifty-fut") {
		enddate = enddate.AddDate(0, 1, 0)
	}

	strEndDate = fmt.Sprintf("%d-%02d-%02d", enddate.Year(), enddate.Month(), enddate.Day())

	symbolFutStr, _ := db.FetchInstrData(order.Instr,
		uint64(order.Entry),
		ts.CtrlParam.TradeSettings.OptionLevel,
		instrumentType,
		strStartDate,
		strEndDate)

	return symbolFutStr
}

// NIFTY21DECFUT
func deriveFuturesName(order data.TradeSignal, ts data.Strategies, selDate time.Time) string {

	var symbolFutStr string = "FAILED"

	Instr := strings.ReplaceAll(order.Instr, "-FUT", "") // remove -FUT suffix

	wkday := selDate.Weekday()
	currThu := time.Now() // dummy initialisation

	if wkday <= time.Thursday {
		currThu = selDate.AddDate(0, ts.CtrlParam.TradeSettings.FuturesExpiryMonth, int(time.Thursday-wkday)) //  upcoming Thu
	} else {
		currThu = selDate.AddDate(0, ts.CtrlParam.TradeSettings.FuturesExpiryMonth, int(7-(wkday-time.Thursday))) //  recent passed Thu + 7 days
	}
	nextThu := currThu.AddDate(0, 0, 7)

	if nextThu.Month() == currThu.Month() { // curr and next thu in same month?

		symbolFutStr = Instr + currThu.Format("06-Jan") + "FUT"

	} else {
		if ts.CtrlParam.TradeSettings.SkipExipryWeekFutures {
			symbolFutStr = Instr + nextThu.Format("06-Jan") + "FUT"
		} else {
			symbolFutStr = Instr + currThu.Format("06-Jan") + "FUT"
		}
	}

	symbolFutStr = strings.ReplaceAll(symbolFutStr, "-", "")
	symbolFutStr = strings.ToUpper(symbolFutStr)

	return symbolFutStr
}
