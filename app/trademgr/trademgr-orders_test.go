package trademgr

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/db"
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"fmt"
	"testing"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

type PlaceOrderTesting struct {
	argStrategy    string
	argInstr       string
	argWeekSel     int
	argStrikePrice float64
	argOptionLevel int
	argDirection   string
	argOrderRoute  string
	argMonthSel    int
	argSkipExpWk   bool
	argOrderType   string
	argVarieties   string
	argValidities  string
	argProducts    string

	// expected          float64
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var PlaceOrderTests = []PlaceOrderTesting{

	// {"2022-03-22", "INFY", 1, 73, 0, "bullish", "equity", 0, false, kiteconnect.VarietyRegular, kiteconnect.ProductMIS},
	// {"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 36000 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "option-buy", 0 + MonthSel, SkipExpWkTrue,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 36000 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "option-sell", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "futures", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "BANKNIFTY-FUT", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bearish", "futures", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "ASHOKLEY-FUT", 0 + WeekSel, 126 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "option-buy", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "ASHOKLEY-FUT", 0 + WeekSel, 126 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "option-sell", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "ASHOK LEYLAND", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductMIS},

	// {"2022-04-02", "ASHOK LEYLAND", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	{"SRB-001", "VODAFONE IDEA", 0 + WeekSel, 10 + StrikePrice, 0 + OptionLevel,
		"bullish", "equity", 0 + MonthSel, SkipExpWkFalse, kiteconnect.OrderTypeLimit,
		kiteconnect.VarietyAMO, kiteconnect.ValidityDay, kiteconnect.ProductMIS},

	// {"RELIANCE INDUSTRIES", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductCNC},

	// {"2022-04-02", "invalid", 0 + WeekSel, 0 + StrikePrice, 0 + OptionLevel,
	// 	"bullish", "equity", 0 + MonthSel, SkipExpWkFalse,
	// 	kiteconnect.VarietyRegular, kiteconnect.ProductCNC},
}

func TestPlaceOrder(t *testing.T) {

	srv.Init()
	srv.LoadEnvVariables("/home/parag/devArea/algotrading-exchange-manager/app/zfiles/config/userSettings.env")
	db.DbInit()
	kite.Init()
	t.Parallel()

	var order appdata.TradeSignal
	var ts appdata.Strategies

	var warCheck = "Check results MANUALLY!!!"

	fmt.Println(ErrorColor, warCheck)
	for _, test := range PlaceOrderTests {

		ts.Strategy = test.argStrategy
		ts.CtrlParam.KiteSettings.Varieties = test.argVarieties
		ts.CtrlParam.KiteSettings.Products = test.argProducts
		ts.CtrlParam.KiteSettings.Validities = test.argValidities
		ts.CtrlParam.KiteSettings.OrderType = test.argOrderType
		ts.CtrlParam.TradeSettings.FuturesExpiryMonth = test.argMonthSel
		ts.CtrlParam.TradeSettings.SkipExipryWeekFutures = test.argSkipExpWk
		ts.CtrlParam.TradeSettings.OrderRoute = test.argOrderRoute
		ts.CtrlParam.TradeSettings.OptionExpiryWeek = test.argWeekSel
		ts.CtrlParam.TradeSettings.OptionLevel = test.argOptionLevel

		order.Dir = test.argDirection
		order.Instr = test.argInstr
		order.Entry = test.argStrikePrice

		// expected := test.expected

		orderID := PlaceOrder(order, ts, time.Now())

		if orderID == 0 {
			t.Errorf(ErrorColor, "\nderiveFuturesName() No data fetched - check dates and levels are correct. This UT is live with server\n")
		} else {
			// print result for manual check
			fmt.Println()
			fmt.Printf(InfoColor, order.Instr)
			fmt.Printf(InfoColor, test.argOrderRoute)
			fmt.Printf(InfoColorUint, orderID)
		}

	}
	fmt.Println()
}
