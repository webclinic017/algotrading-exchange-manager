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
	err := ticker.Subscribe([]uint32{273929, 274185})
	if err != nil {
		fmt.Println("err: ", err)
	}
	err = ticker.SetMode("full", []uint32{273929, 274185})
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
