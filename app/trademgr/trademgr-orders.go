package trademgr

import "goTicker/app/data"

func placeOrder(order *data.TradeSignal, ts *data.Strategies) bool {
	// [x] select opt/fut/stk based on intrument
	// [ ] Fetch account balance
	// [ ] calculate margin required
	// [ ] Check strategy winning percentage
	// [ ] Determine order size
	// [ ] place order
	// [ ] update order id into order table
	// [ ] return order id

	println(ts.CtrlParam.Trades.Base)

	return true

}
