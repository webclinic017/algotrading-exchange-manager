package kite

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/srv"
	"fmt"
	"strconv"
	"strings"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var (
	TickerCnt      uint32 = 0
	ticker         *kiteticker.Ticker
	subscribeToken []uint32
	instrMap            = make(map[string]string)
	Status         bool = false
)

// Triggered when any error is raised
func onError(err error) {
	srv.ErrorLogger.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	Status = false
	srv.InfoLogger.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	srv.InfoLogger.Printf("Connected")
	err := ticker.Subscribe(subscribeToken)
	//err := ticker.Subscribe([]uint32{18257666})
	if err != nil {
		srv.ErrorLogger.Println("err: ", err)
	}
	err = ticker.SetMode("full", subscribeToken)
	//err = ticker.SetMode("full", []uint32{18257666})
	if err != nil {
		srv.ErrorLogger.Println("err: ", err)
	}
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {

	TickerCnt++
	instr := instrMap[fmt.Sprint(tick.InstrumentToken)]
	// fmt.Println("Tick: ", tick)

	if strings.Contains(instr, "-FUT") {
		appdata.ChNseTicks <- appdata.TickData{
			Timestamp:       tick.Timestamp.Time,
			Symbol:          instr,
			LastTradedPrice: tick.LastPrice,
			Buy_Demand:      tick.TotalBuyQuantity,
			Sell_Demand:     tick.TotalSellQuantity,
			TradesTillNow:   tick.VolumeTraded,
			OpenInterest:    tick.OI}
	} else {

		appdata.ChStkTick <- appdata.TickData{
			Timestamp:       tick.Timestamp.Time,
			Symbol:          instr,
			LastTradedPrice: tick.LastPrice,
			Buy_Demand:      tick.TotalBuyQuantity,
			Sell_Demand:     tick.TotalSellQuantity,
			TradesTillNow:   tick.VolumeTraded,
			OpenInterest:    tick.OI}
	}
	// srv.InfoLogger.Println("Time: ", tick.Timestamp.Time, "Instrument: ", instrMap[fmt.Sprint(tick.InstrumentToken)], "LastPrice: ", tick.LastPrice, "Open: ", tick.OHLC.Open, "High: ", tick.OHLC.High, "Low: ", tick.OHLC.Low, "Close: ", tick.OHLC.Close)

	// // Total Buy Quantity, Total Sell quantity, Volume traded, Turnover, Open Interest
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
	Status = false

}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	srv.InfoLogger.Println("Order: ", order.OrderID)
}

func TickerInitialize(apiKey, accToken string) {

	// get current tokens
	instrMap = db.GetInstrumentsToken()
	subscribeToken = getTokens(instrMap)

	// Create new Kite ticker instance
	ticker = kiteticker.New(apiKey, accToken)
	appdata.ChNseTicks = make(chan appdata.TickData, 1000)
	appdata.ChStkTick = make(chan appdata.TickData, 1000)

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
	Status = true
	go ticker.Serve()

}

func CloseTicker() bool {
	defer func() {
		if err := recover(); err != nil {
			srv.InfoLogger.Printf("Boss, ERR in termination of ticker, start debugging :-) ")
		}
	}()
	// ticker.SetAutoReconnect(false)
	Status = false

	ticker.Stop()
	time.Sleep(time.Second * 3) // delay for ticker to terminte connection before we close channel
	close(appdata.ChNseTicks)
	close(appdata.ChStkTick)
	srv.InfoLogger.Printf("Ticker closed for the day, hush!!!")

	return false
}

func TestTicker() {

	// for i := 1; i < 3866; i++ {
	// 	appdata.ChTick <- appdata.TickData{
	// 		Timestamp:       time.Now(),
	// 		Symbol:          "TEST_Signal",
	// 		LastTradedPrice: float64(i),
	// 		Buy_Demand:      uint32(3866 - i),
	// 		Sell_Demand:     12,
	// 		TradesTillNow:   13,
	// 		OpenInterest:    14}

	// 	time.Sleep(time.Millisecond * 1)
	// }
	// close(appdata.ChTick)
}

// From instrMap, get the tokens
func getTokens(instrMap map[string]string) []uint32 {

	//TODO: check if duplicate tokens are present
	var tkn []uint32
	var i = 0

	for key, _ := range instrMap {
		val, _ := strconv.ParseUint(key, 10, 64)
		tkn = append(tkn, uint32(val))
		i++
	}
	return tkn
}
