package kite

import (
	"fmt"
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

func PlaceOrder(orderParams kiteconnect.OrderParams, variety string) uint64 {

	orderResponse, poerr := kc.PlaceOrder(variety, orderParams)
	if poerr != nil {
		print(poerr)
		return 0
	}

	// Resp : convert string to uint
	s, err := strconv.ParseUint(orderResponse.OrderID, 10, 64)
	if err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	return s
}

// RULE "TOTP is mandatory to place orders on third-party apps.
