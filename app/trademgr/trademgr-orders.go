package trademgr

import (
	"goTicker/app/data"
	"goTicker/app/kite"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func placeOrder(order *data.TradeSignal, ts *data.Strategies) bool {
	// [x] select opt/fut/stk based on intrument
	// [ ] Fetch account balance
	// [ ] calculate margin required
	// [ ] Check strategy winning percentage
	// [ ] Determine order size
	// [ ] place order
	// [ ] update order id into order table
	// [ ] return order id

	// println(ts.CtrlParam.Trades.Base)

	return true

}

func PlaceOrder(order data.TradeSignal, ts data.Strategies, selDate time.Time) (orderID uint64) {

	var orderParam kiteconnect.OrderParams

	orderParam.Tag = ts.Strategy
	orderParam.Product = ts.CtrlParam.KiteSettings.Products
	orderParam.Validity = ts.CtrlParam.KiteSettings.Validities

	if strings.ToLower(order.Dir) == "bullish" {
		orderParam.TransactionType = "BUY"
	} else {
		orderParam.TransactionType = "SELL"
	}

	switch ts.CtrlParam.TradeSettings.OrderRoute {

	default:
		fallthrough

	case "equity":
		orderParam.Price = order.Entry
		orderParam.Exchange = kiteconnect.ExchangeNSE
		orderParam.OrderType = ts.CtrlParam.KiteSettings.OrderType

	case "option-buy":
		orderParam.TransactionType = "BUY"
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = kiteconnect.OrderTypeMarket

	case "option-sell":
		orderParam.TransactionType = "SELL"
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = kiteconnect.OrderTypeMarket

	case "futures":
		orderParam.Price = order.Entry
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = ts.CtrlParam.KiteSettings.OrderType

	}
	var symbolMinQty float64
	orderParam.Tradingsymbol, symbolMinQty = deriveInstrumentsName(order, ts, time.Now())
	orderParam.Quantity = int(symbolMinQty)

	return kite.PlaceOrder(orderParam, ts.CtrlParam.KiteSettings.Varieties)

}

// RULE - For optons its always MARKET price, else we need to scan the selected symbol and quote price
