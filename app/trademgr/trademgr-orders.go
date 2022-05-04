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

func pendingOrder(order *appdata.OrderBook_S, ts appdata.UserStrategies_S) bool {

	if !order.Info.Order_simulation {
		return true
	} else {

		/*
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
		*/
		return false
	}
}

func tradeEnter(order *appdata.OrderBook_S, us appdata.UserStrategies_S) bool {

	if !order.Info.Order_simulation { // real trade

		entryTime := time.Now()

		userMargin := kite.GetUserMargin()

		orderMargin := getOrderMargin(*order, us, entryTime)

		order.Info.QtyReq = determineOrderSize(userMargin, orderMargin[0].Total,
			us.Parameters.Controls.WinningRatio, us.Parameters.Controls.MaxBudget,
			us.Parameters.Controls.LimitAmount)

		orderId := executeOrder(*order, us, entryTime, order.Info.QtyReq)

		if orderId != 0 {
			order.Info.OrderIdEntr = orderId
			srv.TradesLogger.Print("Order Placed: ", order.Strategy, " ", orderId)
		}
		return orderId != 0
	} else { // simulation

		order.Info.TradingSymbol = order.Instr
		if strings.Contains(order.Instr, "-FUT") { // RULE: only for futures and equity supported
			order.Info.Exchange = kiteconnect.ExchangeNFO
		} else {
			order.Info.Exchange = kiteconnect.ExchangeNSE
		}
		order.Info.OrderIdEntr = 0
		order.Info.OrderIdExit = 0
		order.Info.QtyReq = 0
		order.Info.QtyFilled = 0
		val, n := kite.GetLatestQuote(order.Instr) // TODO: Add logic to loop through lowest values and return only the price. Add for buy sell
		order.Info.AvgPriceEnter = val[n].Depth.Buy[4].Price
		return true
	}
}

func tradeExit(order *appdata.OrderBook_S, ts appdata.UserStrategies_S) bool {

	if !ts.Parameters.Controls.TradeSimulate {

		entryTime := time.Now()

		userMargin := kite.GetUserMargin()

		orderMargin := getOrderMargin(*order, ts, entryTime)

		order.Info.QtyReq = determineOrderSize(userMargin, orderMargin[0].Total,
			ts.Parameters.Controls.WinningRatio, ts.Parameters.Controls.MaxBudget,
			ts.Parameters.Controls.LimitAmount)

		orderId := executeOrder(*order, ts, entryTime, order.Info.QtyReq)

		if orderId != 0 {
			order.Info.OrderIdEntr = orderId
			srv.TradesLogger.Print("Order Placed: ", order.Strategy, " ", orderId)
		}
		return orderId != 0
	} else {
		return true
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

func executeOrder(order appdata.OrderBook_S, ts appdata.UserStrategies_S, selDate time.Time, qty float64) (orderID uint64) {

	var orderParam kiteconnect.OrderParams

	orderParam.Tag = ts.Strategy
	orderParam.Product = ts.Parameters.Kite_Setting.Products
	orderParam.Validity = ts.Parameters.Kite_Setting.Validities

	if strings.ToLower(order.Dir) == "bullish" {
		orderParam.TransactionType = "BUY"
	} else {
		orderParam.TransactionType = "SELL"
	}

	switch ts.Parameters.Option_setting.OrderRoute {

	default:
		fallthrough

	case "equity":
		orderParam.Price = order.Targets.Entry
		orderParam.Exchange = kiteconnect.ExchangeNSE
		orderParam.OrderType = ts.Parameters.Kite_Setting.OrderType

	case "option-buy":
		orderParam.TransactionType = "BUY"
		orderParam.Exchange = kiteconnect.ExchangeNFO
		orderParam.OrderType = kiteconnect.OrderTypeMarket

	case "option-sell":
		orderParam.TransactionType = "SELL"
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

	return kite.ExecOrder(orderParam, ts.Parameters.Kite_Setting.Varieties)

}

// RULE - For optons its always MARKET price, else we need to scan the selected symbol and quote price
