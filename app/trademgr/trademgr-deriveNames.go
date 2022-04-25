package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"fmt"

	"strings"
	"time"
)

// The format is BANKNIFTY<YY><M><DD>strike<PE/CE>
// The month format is 1 for JAN, 2 for FEB, 3, 4, 5, 6, 7, 8, 9, O(capital o) for October, N for November, D for December.
// var symbolFutStr string = "FAILED"
// BANKNIFTY2232435000CE - 24th Mar 2022
// BANKNIFTY22MAR31000CE - 31st Mar 2022
// Last week of Month - will be monthly expiry
func deriveInstrumentsName(order appdata.OrderBook_S, ts appdata.UserStrategies_S, selDate time.Time) (name string, qty float64) {

	var (
		instrumentType string
		strStartDate   string
		strEndDate     string
		enddate        time.Time
	)

	// ----------------------------------------------------------------------
	if ts.Parameters.Option_setting.OrderRoute == "option-buy" {
		selDate = selDate.AddDate(0, 0, (7 * ts.Parameters.Option_setting.OptionExpiryWeek))
		enddate = selDate.AddDate(0, 0, 7+(7*ts.Parameters.Option_setting.OptionExpiryWeek))
		// ---------------------------------------------------------------------- Special case for expiry
		// For individual securities expiry is monthly
		if (strings.ToLower(order.Instr) != "nifty-fut") ||
			(strings.ToLower(order.Instr) != "banknifty-fut") ||
			(strings.ToLower(order.Instr) != "finnifty-fut") ||
			(strings.ToLower(order.Instr) != "midcpnifty-fut") {
			enddate = selDate.AddDate(0, 1, 0)
		}
		if strings.ToLower(order.Dir) == "bullish" {
			instrumentType = "CE"
		} else {
			instrumentType = "PE"
		}
	} else if ts.Parameters.Option_setting.OrderRoute == "option-sell" {
		selDate = selDate.AddDate(0, 0, (7 * ts.Parameters.Option_setting.OptionExpiryWeek))
		enddate = selDate.AddDate(0, 0, 7+(7*ts.Parameters.Option_setting.OptionExpiryWeek))
		// ---------------------------------------------------------------------- Special case for expiry
		// For individual securities expiry is monthly
		if (strings.ToLower(order.Instr) != "nifty-fut") ||
			(strings.ToLower(order.Instr) != "banknifty-fut") ||
			(strings.ToLower(order.Instr) != "finnifty-fut") ||
			(strings.ToLower(order.Instr) != "midcpnifty-fut") {
			enddate = selDate.AddDate(0, 1, 0)
		}
		if strings.ToLower(order.Dir) == "bullish" {
			instrumentType = "PE"
		} else {
			instrumentType = "CE"
		}
	} else if ts.Parameters.Option_setting.OrderRoute == "futures" {
		selDate = selDate.AddDate(0, ts.Parameters.Futures_Setting.FuturesExpiryMonth, 0)
		enddate = selDate.AddDate(0, 1+ts.Parameters.Futures_Setting.FuturesExpiryMonth, 0)
		instrumentType = "FUT"
	} else if ts.Parameters.Option_setting.OrderRoute == "equity" {
		enddate = selDate.AddDate(0, 0, 0)
		instrumentType = "EQ"
	}

	strStartDate = fmt.Sprintf("%d-%02d-%02d", selDate.Year(), selDate.Month(), selDate.Day())

	strEndDate = fmt.Sprintf("%d-%02d-%02d", enddate.Year(), enddate.Month(), enddate.Day())

	symbolFutStr, qty := db.FetchInstrData(order.Instr,
		uint64(order.Targets.Entry),
		ts.Parameters.Option_setting.OptionLevel,
		instrumentType,
		strStartDate,
		strEndDate)

	return symbolFutStr, qty
}
