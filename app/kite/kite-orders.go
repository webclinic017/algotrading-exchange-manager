package kite

import (
	"algo-ex-mgr/app/srv"
	"strconv"
	"strings"

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

func ModifyOrder(orderId uint64, variety string, orderParams kiteconnect.OrderParams) uint64 {
	// TODO; testing pending
	odr := strconv.FormatUint(orderId, 10)
	orderResponse, err := kc.ModifyOrder(variety, odr, orderParams)
	if err != nil {
		srv.TradesLogger.Println("ModifyOrder: ", err.Error())
		return 0
	}
	s, err := strconv.ParseUint(orderResponse.OrderID, 10, 64)
	if err == nil {
		srv.TradesLogger.Printf("%T, %v\n", s, s)
	}
	return s

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

func GetLatestQuote(i string) (kiteconnect.Quote, string) {

	var q string
	if strings.Contains(i, "-FUT") {
		q = "NFO:" + strings.Replace(i, "-FUT", "", -1)
	} else {
		q = "NSE:" + i
	}

	quote, err := kc.GetQuote(q) // 'exchange:Insturment'
	if err != nil {
		srv.TradesLogger.Println(err.Error())
		return kiteconnect.Quote{}, q
	}
	return quote, q
}
