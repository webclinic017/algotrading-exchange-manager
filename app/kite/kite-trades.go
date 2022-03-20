package kite

import (
	"goTicker/app/data"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func PlaceOrder(order *data.TradeSignal) bool {

	if order.Instr != "" {
		return true
	} else {
		return false
	}
}

func CalOrderMargin(order data.TradeSignal, ts data.Strategies) bool {

	var marginParam kiteconnect.GetMarginParams

	// default params
	marginParam.Compact = false
	marginParam.OrderParams[0].Exchange = "NSE"
	marginParam.OrderParams[0].OrderType = "MARKET"
	marginParam.OrderParams[0].Quantity = 1
	marginParam.OrderParams[0].Price = 0
	marginParam.OrderParams[0].TriggerPrice = 0
	// specific params
	marginParam.OrderParams[0].Variety = ts.CtrlParam.KiteSettings.Varieties
	marginParam.OrderParams[0].Product = ts.CtrlParam.KiteSettings.Products
	marginParam.OrderParams[0].TransactionType = order.Dir

	switch ts.CtrlParam.TradeSettings.OrderRoute {

	default:
		fallthrough
	case "stock":
		marginParam.OrderParams[0].Tradingsymbol = order.Instr

	case "option":
		marginParam.OrderParams[0].Tradingsymbol = deriveOptionName(order, ts, time.Now())

	case "futures":
		marginParam.OrderParams[0].Tradingsymbol = deriveFuturesName(order, ts, time.Now())

	}
	OrderMargins, err := kc.GetOrderMargins(marginParam)

	print(OrderMargins, err)
	return true

}

func deriveOptionName(order data.TradeSignal, ts data.Strategies, selDate time.Time) string {
	// The format is BANKNIFTY<YY><M><DD>strike<PE/CE>
	// The month format is 1 for JAN, 2 for FEB, 3, 4, 5, 6, 7, 8, 9, O(capital o) for October, N for November, D for December.

	return ""
}

func deriveFuturesName(order data.TradeSignal, ts data.Strategies, selDate time.Time) string {

	var symbolFutStr string = "FAILED"
	// NIFTY21DECFUT

	Instr := strings.ReplaceAll(order.Instr, "-FUT", "") // remove -FUT suffix
	// monthSelected := selDate.AddDate(0, ts.CtrlParam.TradeSettings.FuturesExpiryMonth, 0)

	// currThu := selDate.AddDate(0, ts.CtrlParam.TradeSettings.FuturesExpiryMonth, 0)

	wkday := selDate.Weekday()
	currThu := time.Now()

	if wkday <= time.Thursday {
		currThu = selDate.AddDate(0, ts.CtrlParam.TradeSettings.FuturesExpiryMonth, int(time.Thursday-wkday)) //  upcoming Thu + 7 days
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
