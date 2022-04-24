package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func pendingOrder(order *appdata.OrderBook_S, ts appdata.UserStrategies_S) bool {

	tradesList := kite.FetchOrderTrades(order.Order_id)
	var qtyFilled float64

	for each := range tradesList {
		qtyFilled = qtyFilled + tradesList[each].Quantity
	}
	order.OrderData.QtyFilled = qtyFilled
	od, err := json.Marshal(tradesList)
	if err != nil {
		srv.TradesLogger.Println("Error in marshalling trades list: ", err.Error())
		return false
	}
	order.Order_trades_entry = string(od)

	if order.OrderData.QtyReq > order.OrderData.QtyFilled {
		// TODO: modify limit price if order is still pending
		qt, n := kite.GetLatestQuote(order.Instr)
		// ToDO: fetch min values for the enter/exit trades
		fmt.Print(qt[n].Depth.Buy[0].Price)

		return true
	} else {
		return false
	}
}

func tradeEnter(order *appdata.OrderBook_S, ts appdata.UserStrategies_S) bool {

	entryTime := time.Now()

	userMargin := kite.GetUserMargin()

	orderMargin := getOrderMargin(*order, ts, entryTime)

	order.OrderData.QtyReq = determineOrderSize(userMargin, orderMargin[0].Total,
		ts.CtrlData.Percentages.WinningRate, ts.CtrlData.Percentages.MaxBudget,
		ts.CtrlData.Trade_Setting.LimitAmount)

	orderId := executeOrder(*order, ts, entryTime, order.OrderData.QtyReq)

	if orderId != 0 {
		order.Order_id = orderId
		srv.TradesLogger.Print("Order Placed: ", order.Strategy, " ", orderId)
	}
	return orderId != 0
}

// Fetch account balance
// Calculate margin required
// Check strategy winning percentage
// Determine order size
func determineOrderSize(userMargin float64, orderMargin float64, winningRate float64, maxBudget float64, limitAmount float64) float64 {

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
				return qty // based on winning rate
			}
		}
	}
}

func executeOrder(order appdata.OrderBook_S, ts appdata.UserStrategies_S, selDate time.Time, qty float64) (orderID uint64) {

	var orderParam kiteconnect.OrderParams

	orderParam.Tag = ts.Strategy
	orderParam.Product = ts.CtrlData.Kite_Setting.Products
	orderParam.Validity = ts.CtrlData.Kite_Setting.Validities

	if strings.ToLower(order.Dir) == "bullish" {
		orderParam.TransactionType = "BUY"
	} else {
		orderParam.TransactionType = "SELL"
	}

	switch ts.CtrlData.Trade_Setting.OrderRoute {

	default:
		fallthrough

	case "equity":
		orderParam.Price = order.Entry
		orderParam.Exchange = kiteconnect.ExchangeNSE
		orderParam.OrderType = ts.CtrlData.Kite_Setting.OrderType

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
		orderParam.OrderType = ts.CtrlData.Kite_Setting.OrderType

	}
	var symbolMinQty float64
	orderParam.Tradingsymbol, symbolMinQty = deriveInstrumentsName(order, ts, time.Now())
	orderParam.Quantity = int(symbolMinQty * qty)

	return kite.ExecOrder(orderParam, ts.CtrlData.Kite_Setting.Varieties)

}

// RULE - For optons its always MARKET price, else we need to scan the selected symbol and quote price
