package kite

import (
	"fmt"
	"goTicker/app/srv"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var (
	TickerCnt            uint32 = 0
	ticker               *kiteticker.Ticker
	Tokens               []uint32
	TokensWithNames      []string
	ChTick               chan TickData
	InsNamesMap          = make(map[string]string)
	symbolFutStr         string
	symbolMcxFutStr      string
	KiteConnectionStatus bool = false
)

type TickData struct {
	Timestamp       time.Time
	LastTradedPrice float64
	Symbol          string
	LastPrice       float64
	Buy_Demand      uint32
	Sell_Demand     uint32
	TradesTillNow   uint32
	OpenInterest    uint32
}

// Triggered when any error is raised
func onError(err error) {
	srv.ErrorLogger.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	KiteConnectionStatus = false
	srv.InfoLogger.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	srv.InfoLogger.Printf("Connected")
	err := ticker.Subscribe(Tokens)
	//err := ticker.Subscribe([]uint32{18257666})
	if err != nil {
		srv.ErrorLogger.Println("err: ", err)
	}
	err = ticker.SetMode("full", Tokens)
	//err = ticker.SetMode("full", []uint32{18257666})
	if err != nil {
		srv.ErrorLogger.Println("err: ", err)
	}
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {

	TickerCnt++

	ChTick <- TickData{
		Timestamp:       tick.Timestamp.Time,
		Symbol:          InsNamesMap[fmt.Sprint(tick.InstrumentToken)],
		LastTradedPrice: tick.LastPrice,
		Buy_Demand:      tick.TotalBuyQuantity,
		Sell_Demand:     tick.TotalSellQuantity,
		TradesTillNow:   tick.VolumeTraded,
		OpenInterest:    tick.OI}

	// srv.InfoLogger.Println("Time: ", tick.Timestamp.Time, "Instrument: ", InsNamesMap[fmt.Sprint(tick.InstrumentToken)], "LastPrice: ", tick.LastPrice, "Open: ", tick.OHLC.Open, "High: ", tick.OHLC.High, "Low: ", tick.OHLC.Low, "Close: ", tick.OHLC.Close)

	// Total Buy Quantity, Total Sell quantity, Volume traded, Turnover, Open Interest
	// fmt.Println("Total Buy Quantity: ", tick.TotalBuyQuantity)
	// fmt.Println("Total Sell Quantity: ", tick.TotalSellQuantity)
	// fmt.Println("VolumeTraded: ", tick.VolumeTraded)
	// fmt.Println("LastTradedQuantity: ", tick.LastTradedQuantity)
	// fmt.Println("TotalBuy: ", tick.TotalBuy)
	// fmt.Println("TotalSell: ", tick.TotalSell)
	// fmt.Println("OI: ", tick.OI)

}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	srv.InfoLogger.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	srv.InfoLogger.Printf("Maximum no of reconnect attempt reached: %d", attempt)
	KiteConnectionStatus = false

}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	srv.InfoLogger.Println("Order: ", order.OrderID)
}

func TickerInitialize(apiKey, accToken string) {

	// Create new Kite ticker instance
	ticker = kiteticker.New(apiKey, accToken)
	ChTick = make(chan TickData, 1000)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	srv.InfoLogger.Printf("Initiaing Ticker connection, time to make money --->>> drama unfurls now...")
	KiteConnectionStatus = true
	go ticker.Serve()

}

func CloseTicker() bool {
	defer func() {
		if err := recover(); err != nil {
			srv.InfoLogger.Printf("Boss, ERR in termination of ticker, start debugging :-) ")
		}
	}()
	// ticker.SetAutoReconnect(false)
	KiteConnectionStatus = false

	ticker.Stop()
	time.Sleep(time.Second * 3) // delay for ticker to terminte connection before we close channel
	close(ChTick)
	srv.InfoLogger.Printf("Ticker closed for the day, hush!!!")

	return false
}

func TestTicker() {

	for i := 1; i < 3866; i++ {
		ChTick <- TickData{
			Timestamp:       time.Now(),
			Symbol:          "TEST_Signal",
			LastTradedPrice: float64(i),
			Buy_Demand:      uint32(3866 - i),
			Sell_Demand:     12,
			TradesTillNow:   13,
			OpenInterest:    14}

		time.Sleep(time.Millisecond * 1)
	}
	close(ChTick)
}
