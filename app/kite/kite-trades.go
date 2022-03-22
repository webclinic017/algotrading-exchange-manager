package kite

import (
	"goTicker/app/data"
)

func PlaceOrder(order *data.TradeSignal) bool {

	if order.Instr != "" {
		return true
	} else {
		return false
	}
}
