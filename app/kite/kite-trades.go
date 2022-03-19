package kite

import (
	"goTicker/app/data"
	"goTicker/app/srv"
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

func CalOrderMargin(order *data.TradeSignal, ts *data.Strategies) bool {

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
		marginParam.OrderParams[0].Tradingsymbol = deriveOptionName(order, ts)

	case "futures":
		marginParam.OrderParams[0].Tradingsymbol = deriveFuturesName(order, ts)

	}
	OrderMargins, err := kc.GetOrderMargins(marginParam)

	print(OrderMargins, err)
	return true

}

func deriveOptionName(order *data.TradeSignal, ts *data.Strategies) string {
	// The format is BANKNIFTY<YY><M><DD>strike<PE/CE>
	// The month format is 1 for JAN, 2 for FEB, 3, 4, 5, 6, 7, 8, 9, O(capital o) for October, N for November, D for December.

	return ""
}

func deriveFuturesName(order *data.TradeSignal, ts *data.Strategies) string {

	var symbolFutStr string = "FAILED"
	// NIFTY21DECFUT

	monthSelected := time.Now().AddDate(0, ts.CtrlParam.TradeSettings.FuturesExpiryMonth, 0)

	wkday := time.Now().Weekday()
	if wkday <= time.Thursday {
		nextThu := time.Now().AddDate(0, 0, int((time.Thursday-wkday)+7)) //  curr Thu + 7 days
		if nextThu.Month() == monthSelected.Month() {                     // curr and next thu in same month?
			symbolFutStr = monthSelected.Format("06-Jan") + "FUT"
		} else {
			if ts.CtrlParam.TradeSettings.SkipExipryWeekFutures {
				symbolFutStr = nextThu.Format("06-Jan") + "FUT"

			} else {
				symbolFutStr = monthSelected.Format("06-Jan") + "FUT"
			}
		}
	}

	symbolFutStr = strings.ReplaceAll(symbolFutStr, "-", "")
	symbolFutStr = strings.ToUpper(symbolFutStr)
	srv.InfoLogger.Println("Futures Symbol : Decoded :- ", symbolFutStr)

	return symbolFutStr
}
