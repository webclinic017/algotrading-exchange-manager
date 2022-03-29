package kite

import (
	"fmt"
	"strconv"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

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

// RULE "TOTP is mandatory to place orders on third-party apps. Learn how to set up TOTP for your account [here](https://support.zerodha.com/category/your-zerodha-account/login-credentials/login-credentials-of-trading-platforms/articles/time-based-otp-setup-mobile-only)."
