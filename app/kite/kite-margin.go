package kite

import (
	"goTicker/app/data"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

// curl https://api.kite.trade/margins/orders \
//     -H 'X-Kite-Version: 3' \
//     -H 'Authorization: token api_key:access_token' \
//     -H 'Content-Type: application/json' \
//     -d '[
//     {
//         "exchange": "NSE",
//         "tradingsymbol": "INFY",
//         "transaction_type": "BUY",
//         "variety": "regular",
//         "product": "CNC",
//         "order_type": "MARKET",
//         "quantity": 1,
//         "price": 0,
//         "trigger_price": 0
//     }
// ]'

func CalOrderMargin(order data.TradeSignal, ts data.Strategies, tm time.Time) []kiteconnect.OrderMargins {

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
	if strings.ToLower(order.Dir) == "bullish" {
		marginParam.OrderParams[0].TransactionType = "BUY"
	} else {
		marginParam.OrderParams[0].TransactionType = "SELL"
	}

	switch ts.CtrlParam.TradeSettings.OrderRoute {

	default:
		fallthrough

	case "stock":
		marginParam.OrderParams[0].Tradingsymbol = order.Instr

	case "option-buy":
		marginParam.OrderParams[0].TransactionType = "BUY"
		marginParam.OrderParams[0].Tradingsymbol = deriveOptionName(order, ts, tm)

	case "option-sell":
		marginParam.OrderParams[0].TransactionType = "SELL"
		marginParam.OrderParams[0].Tradingsymbol = deriveOptionName(order, ts, tm)

	case "futures":
		marginParam.OrderParams[0].Tradingsymbol = deriveFuturesName(order, ts, tm)

	}
	OrderMargins, err := kc.GetOrderMargins(marginParam)

	print(OrderMargins, err)
	return OrderMargins

}
