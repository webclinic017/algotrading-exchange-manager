package kite

import (
	"fmt"
	"log"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var (
	ticker *kiteticker.Ticker

	ChTick = make(chan TickData, 3)
)

type TickData struct {
	Timestamp time.Time
	Lastprice float64
	Insttoken uint32
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    uint32
}

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
	err := ticker.Subscribe([]uint32{273929, 274185, 274441, 274697, 274953, 275209, 275465, 275721, 275977, 276233, 276489, 276745, 277001, 277257, 277513, 277769, 278025, 278281, 278537, 278793, 279049, 279305, 279561, 279817, 280073, 280329, 280585, 280841, 281097, 281353, 281609, 281865, 282121, 282377, 282633, 282889, 283145, 283401, 283657, 283913, 284169, 284425, 284681, 284937, 285193, 285449, 265, 285961, 286217, 286473, 286729, 286985, 287241, 287497, 287753, 273417, 264713, 264969, 260617, 264457, 256265, 268041, 265993, 273673, 263433, 260105, 257289, 257545, 268297, 257033, 261641, 257801, 261897, 270345, 269065, 269321, 269577, 269833, 268553, 268809, 270089, 261385, 259849, 263945, 263689, 270601, 256777, 266249, 260873, 266505, 262153, 270857, 262409, 262665, 262921, 271113, 261129, 263177, 267017, 267273, 266761, 271881, 267785, 272137, 272393, 265737, 265225, 271625, 259081, 258825, 259593, 259337, 267529})
	if err != nil {
		fmt.Println("err: ", err)
	}
	err = ticker.SetMode("full", []uint32{273929, 274185, 274441, 274697, 274953, 275209, 275465, 275721, 275977, 276233, 276489, 276745, 277001, 277257, 277513, 277769, 278025, 278281, 278537, 278793, 279049, 279305, 279561, 279817, 280073, 280329, 280585, 280841, 281097, 281353, 281609, 281865, 282121, 282377, 282633, 282889, 283145, 283401, 283657, 283913, 284169, 284425, 284681, 284937, 285193, 285449, 265, 285961, 286217, 286473, 286729, 286985, 287241, 287497, 287753, 273417, 264713, 264969, 260617, 264457, 256265, 268041, 265993, 273673, 263433, 260105, 257289, 257545, 268297, 257033, 261641, 257801, 261897, 270345, 269065, 269321, 269577, 269833, 268553, 268809, 270089, 261385, 259849, 263945, 263689, 270601, 256777, 266249, 260873, 266505, 262153, 270857, 262409, 262665, 262921, 271113, 261129, 263177, 267017, 267273, 266761, 271881, 267785, 272137, 272393, 265737, 265225, 271625, 259081, 258825, 259593, 259337, 267529})
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {
	//fmt.Println("Tick: ", tick)
	ChTick <- TickData{
		Timestamp: tick.Timestamp.Time,
		Insttoken: tick.InstrumentToken,
		Lastprice: tick.LastPrice,
		Open:      tick.OHLC.Open,
		High:      tick.OHLC.High,
		Low:       tick.OHLC.Low,
		Close:     tick.OHLC.Close,
		Volume:    tick.VolumeTraded}

	/*
		fmt.Println("Time: ", tick.Timestamp.Time)
		fmt.Println("Instrument: ", tick.InstrumentToken)
		fmt.Println("LastPrice: ", tick.LastPrice)
		fmt.Println("Open: ", tick.OHLC.Open)
		fmt.Println("High: ", tick.OHLC.High)
		fmt.Println("Low: ", tick.OHLC.Low)
		fmt.Println("Close: ", tick.OHLC.Close)
		fmt.Println("Volumne: ", tick.VolumeTraded)
	*/

}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Println("Order: ", order.OrderID)
}

func TickerInitialize(apiKey, accToken string) {

	// Create new Kite ticker instance
	ticker = kiteticker.New(apiKey, accToken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	go ticker.Serve()

}

func CloseTicker() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Terminating ticker:", err)
		}
	}()
	// ticker.SetAutoReconnect(false)
	ticker.Stop()
	println("Ticker connection closed")
}
