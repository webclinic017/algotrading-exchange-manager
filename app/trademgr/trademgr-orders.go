package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"math"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func tradeEnter(order appdata.TradeSignal, ts appdata.Strategies) bool {

	entryTime := time.Now()

	userMargin := kite.GetUserMargin()

	orderMargin := getOrderMargin(order, ts, entryTime)

	tradeQty := determineOrderSize(userMargin, orderMargin[0].Total,
		ts.CtrlParam.Percentages.WinningRate, ts.CtrlParam.Percentages.MaxBudget,
		ts.CtrlParam.TradeSettings.LimitAmount)

	orderId := executeOrder(order, ts, entryTime, tradeQty)
	TradesList := kite.FetchOrderTrades(orderId)

	srv.TradesLogger.Print("Trade executed: ", TradesList)

	// println("Order Placed : ", orderId)

	// [ ] update order id into order table
	// [ ] return order id

	// println(ts.CtrlParam.Trades.Base)

	return true

}

// Fetch account balance
// Calculate margin required
// Check strategy winning percentage
// Determine order size
func determineOrderSize(userMargin float64, orderMargin float64, winningRate float64, maxBudget float64, limitAmount float64) int {

	maxBudget = (maxBudget / 100) * userMargin
	budget := math.Min(maxBudget, limitAmount)

	if orderMargin > budget { // no money available for transaction
		return 0
	} else {
		qty := (budget / orderMargin) * (winningRate / 100) // place order in % of winning rate
		if qty < 1 {
			return 1 // minimum order size if winning rate is less than 1
		} else {
			if math.IsNaN(qty) {
				return 0
			} else {
				return int(qty) // based on winning rate
			}
		}
	}
}

func executeOrder(order appdata.TradeSignal, ts appdata.Strategies, selDate time.Time, qty int) (orderID uint64) {

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
	orderParam.Quantity = int(symbolMinQty) * qty

	return kite.ExecOrder(orderParam, ts.CtrlParam.KiteSettings.Varieties)

}

// RULE - For optons its always MARKET price, else we need to scan the selected symbol and quote price
