package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func getOrderMargin(order appdata.OrderBook_S, ts appdata.UserStrategies_S, tm time.Time) []kiteconnect.OrderMargins {

	var marginParam kiteconnect.GetMarginParams

	//initialise the slice
	marginParam.OrderParams = make([]kiteconnect.OrderMarginParam, 1)

	// default params
	marginParam.Compact = false

	marginParam.OrderParams[0].OrderType = "MARKET"
	marginParam.OrderParams[0].Price = 0
	marginParam.OrderParams[0].TriggerPrice = 0
	// specific params
	marginParam.OrderParams[0].Variety = ts.Parameters.Kite_Setting.Varieties
	marginParam.OrderParams[0].Product = ts.Parameters.Kite_Setting.Products
	if strings.ToLower(order.Dir) == "bullish" {
		marginParam.OrderParams[0].TransactionType = "BUY"
	} else {
		marginParam.OrderParams[0].TransactionType = "SELL"
	}

	switch ts.Parameters.Kite_Setting.OrderRoute {

	case "equity":
		marginParam.OrderParams[0].Exchange = kiteconnect.ExchangeNSE

	case "option-buy":
		marginParam.OrderParams[0].Exchange = kiteconnect.ExchangeNFO
		marginParam.OrderParams[0].TransactionType = "BUY"

	case "option-sell":
		marginParam.OrderParams[0].Exchange = kiteconnect.ExchangeNFO
		marginParam.OrderParams[0].TransactionType = "SELL"

	case "futures":
		marginParam.OrderParams[0].Exchange = kiteconnect.ExchangeNFO

	}
	marginParam.OrderParams[0].Tradingsymbol, marginParam.OrderParams[0].Quantity =
		deriveInstrumentsName(order, ts, tm)

	OrderMargins, err := kite.FetchOrderMargins(marginParam)

	if err != nil {
		srv.ErrorLogger.Println(err)
	}
	return OrderMargins

}
