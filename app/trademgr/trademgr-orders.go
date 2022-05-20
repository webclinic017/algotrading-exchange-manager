package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"math"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func pendingOrderEntr(order *appdata.OrderBook_S, us appdata.UserStrategies_S) bool {

	if order.Info.Order_simulation {
		return true
	} else {

		tradesList := kite.FetchOrderTrades(order.Info.OrderIdEntr)
		var qtyFilled float64

		for each := range tradesList {
			qtyFilled = qtyFilled + tradesList[each].Quantity
		}
		order.Info.QtyFilledEntr = qtyFilled

		order.Orders_entr = make([]kiteconnect.Trade, len(tradesList))
		print(copy(order.Orders_entr, tradesList))

		if order.Info.QtyReq > order.Info.QtyFilledEntr {
			_ = finalizeOrder(*order, us, time.Now(), (order.Info.QtyReq - order.Info.QtyFilledEntr), true)
			return false
		} else {
			return true
		}

	}
}

func pendingOrderExit(order *appdata.OrderBook_S, us appdata.UserStrategies_S) bool {

	if order.Info.Order_simulation {
		return true
	} else {

		tradesList := kite.FetchOrderTrades(uint64(order.Info.OrderIdExit))
		var qtyFilled float64

		for each := range tradesList {
			qtyFilled = qtyFilled + tradesList[each].Quantity
		}
		order.Info.QtyFilledExit = qtyFilled

		order.Orders_exit = make([]kiteconnect.Trade, len(tradesList))
		print(copy(order.Orders_exit, tradesList))

		if order.Info.QtyFilledEntr > order.Info.QtyFilledExit {
			_ = finalizeOrder(*order, us, time.Now(), (order.Info.QtyFilledEntr - order.Info.QtyFilledExit), false)
			return false
		} else {
			return true
		}

	}
}

func tradeEnter(order *appdata.OrderBook_S, us appdata.UserStrategies_S) bool {

	if order.Info.Order_simulation { // real trade

		order.Info.TradingSymbol = order.Instr
		if strings.Contains(order.Instr, "-FUT") { // RULE: only for futures and equity supported
			order.Info.Exchange = kiteconnect.ExchangeNFO
		} else {
			order.Info.Exchange = kiteconnect.ExchangeNSE
		}
		order.Info.OrderIdEntr = 0
		order.Info.OrderIdExit = 0
		order.Info.QtyReq = 0
		order.Info.QtyFilledEntr = 0
		val, n := kite.GetLatestQuote(order.Instr) // TODO: Add logic to loop through lowest values and return only the price. Add for buy sell
		order.Info.AvgPriceEnter = val[n].Depth.Buy[4].Price
		return true

	} else { // simulation
		entryTime := time.Now()

		userMargin := kite.GetUserMargin()
		orderMargin := getOrderMargin(*order, us, entryTime)

		order.Info.QtyReq = determineOrderSize(userMargin, orderMargin[0].Total,
			us.Parameters.Controls.WinningRatio, us.Parameters.Controls.MaxBudget,
			us.Parameters.Controls.LimitAmount)

		orderId := finalizeOrder(*order, us, entryTime, order.Info.QtyReq, true)

		if orderId != 0 {
			order.Info.OrderIdEntr = orderId
			srv.TradesLogger.Print("Order Placed: ", order.Strategy, " ", orderId)
		}
		return orderId != 0
	}
}

func tradeExit(order *appdata.OrderBook_S, ts appdata.UserStrategies_S) bool {

	if ts.Parameters.Controls.TradeSimulate {
		return true
	} else {

		orderId := finalizeOrder(*order, ts, time.Now(), order.Info.QtyFilledEntr, false)

		if orderId != 0 {
			order.Info.OrderIdExit = orderId
			srv.TradesLogger.Print("Order Placed: ", order.Strategy, " ", orderId)
		}
		return orderId != 0
	}
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
				return math.Trunc(qty) // based on winning rate
			}
		}
	}
}

func getLowestPrice(instr string) float64 {
	qt, n := kite.GetLatestQuote(instr)
	// ToDO: fetch min values for the enter/exit trades
	fmt.Print(qt[n].Depth.Buy[0].Price)
	return qt[n].Depth.Buy[0].Price
}

func finalizeOrder(order appdata.OrderBook_S, ts appdata.UserStrategies_S, selDate time.Time, qty float64, entry bool) (orderID uint64) {

	var orderParam kiteconnect.OrderParams

	orderParam.Tag = ts.Strategy
	orderParam.Product = ts.Parameters.Kite_Setting.Products
	orderParam.Validity = ts.Parameters.Kite_Setting.Validities

	// Valid only for equity and futures
	// entry(true) + bullish - buy
	// exit(false) + bearish - buy
	// exit(false) + bullish - sell
	// entry(true) + bearish - sell

	if (strings.ToLower(order.Dir) == "bullish" && !entry) || (strings.ToLower(order.Dir) == "bearish" && entry) {
		orderParam.TransactionType = "SELL"
	} else {
		orderParam.TransactionType = "BUY"
	}

	switch ts.Parameters.Kite_Setting.OrderRoute {

	default:
		fallthrough

	case "equity":
		orderParam.Price = order.Targets.Entry
		orderParam.Exchange = kiteconnect.ExchangeNSE
		orderParam.OrderType = ts.Parameters.Kite_Setting.OrderType

	case "option-buy":
		if entry {
			orderParam.TransactionType = "BUY"
		} else {
			orderParam.TransactionType = "SELL"
		}
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = kiteconnect.OrderTypeMarket

	case "option-sell":
		if entry {
			orderParam.TransactionType = "SELL"
		} else {
			orderParam.TransactionType = "BUY"
		}
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = kiteconnect.OrderTypeMarket

	case "futures":
		orderParam.Price = order.Targets.Entry
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = ts.Parameters.Kite_Setting.OrderType

	}
	var symbolMinQty float64
	orderParam.Tradingsymbol, symbolMinQty = deriveInstrumentsName(order, ts, time.Now())
	orderParam.Quantity = int(symbolMinQty * qty)

	var orderId uint64
	if entry { // check if order is already placed
		orderId = order.Info.OrderIdEntr
	} else {
		orderId = order.Info.OrderIdExit
	}

	if orderId == 0 { // new order
		return kite.ExecOrder(orderParam, ts.Parameters.Kite_Setting.Varieties)
	} else {
		return kite.ModifyOrder(orderId, ts.Parameters.Kite_Setting.Varieties, orderParam)
	}
}

// RULE - For optons its always MARKET price, else we need to scan the selected "option symbol" and quote price
