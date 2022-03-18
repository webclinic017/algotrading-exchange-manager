package kite

import (
	"goTicker/app/data"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func PlaceOrder(order *data.TradeSignal) bool {

	if order.Instr != "" {
		return true
	} else {
		return false
	}
}

func CalOrderMargin(order *data.TradeSignal) bool {

	var marginParam kiteconnect.GetMarginParams

	if order.Instr != "" {
		return true
	} else {
		return false
	}
}
