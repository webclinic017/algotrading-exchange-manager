package kite

import (
	"fmt"
	"strconv"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func PlaceOrder(orderParams kiteconnect.OrderParams) uint64 {

	orderResponse, poerr := kc.PlaceOrder(kiteconnect.VarietyRegular, orderParams)
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
