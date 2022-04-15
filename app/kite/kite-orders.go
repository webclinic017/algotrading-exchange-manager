package kite

import (
	"algo-ex-mgr/app/srv"
	"strconv"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func FetchOrderMargins(marginParam kiteconnect.GetMarginParams) ([]kiteconnect.OrderMargins, error) {
	var OrderMargins []kiteconnect.OrderMargins
	var err error

	OrderMargins, err = kc.GetOrderMargins(marginParam)

	if err != nil {
		return nil, err
	}

	return OrderMargins, nil
}

func FetchOrderTrades(orderId uint64) []kiteconnect.Trade {
	odr := strconv.FormatUint(orderId, 10)
	order, err := kc.GetOrderTrades(odr)
	if err != nil {
		srv.TradesLogger.Println(err.Error())
		return nil
	}
	srv.TradesLogger.Println(order)
	return order
}

func ExecOrder(orderParams kiteconnect.OrderParams, variety string) uint64 {

	orderResponse, poerr := kc.PlaceOrder(variety, orderParams)
	if poerr != nil {
		srv.TradesLogger.Println(poerr.Error())
		return 0
	}

	// Resp : convert string to uint
	s, err := strconv.ParseUint(orderResponse.OrderID, 10, 64)
	if err == nil {
		srv.TradesLogger.Printf("%T, %v\n", s, s)
	}
	return s
}

// RULE "TOTP is mandatory to place orders on third-party apps.
