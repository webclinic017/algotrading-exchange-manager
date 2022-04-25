package appdata

import (
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

// Global variables
const (
	InfoColorFloat = "\033[1;34m%.0f\033[0m\t"
	InfoColorUint  = "\033[1;34m%d\033[0m\t"
	InfoColor      = "\033[1;34m%20s\033[0m\t"
	SuccessColor   = "\033[1;32m%20s\033[0m\t"
	ErrorColor     = "\033[1;31m%s\033[0m"
	DebugColor     = "\033[1;35m%s\033[0m"

	ColorSuccess = "\033[37m\033[42m"
	ColorBanner  = "\033[34;4m"
	ColorInfo    = "\033[37m"
	ColorDimmed  = "\033[37m"
	ColorReset   = "\033[0m"
	ColorError   = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorWarning = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorPurple  = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m\033[47m"
)

// Foreground text colors
// const (
// 	FgBlack Attribute = iota + 30
// 	FgRed
// 	FgGreen
// 	FgYellow
// 	FgBlue
// 	FgMagenta
// 	FgCyan
// 	FgWhite
// )

var (
	ChNseTicks chan TickData
	ChStkTick  chan TickData
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

type Percentage_S struct {
	Target      float64 `json:"Target"`      // "target": 1,
	SL          float64 `json:"SL"`          // "sl": 1,
	DeepSL      float64 `json:"DeepSL"`      // "deepsl": 1
	MaxBudget   float64 `json:"MaxBudget"`   // "limit_budget": 50%,
	WinningRate float64 `json:"WinningRate"` // "winning_rate": 80%,
}

type TargetControls_S struct {
	Trail_target_en         bool      `json:"Trail_target_en"`         // 	"trail_target_en": true,
	Position_reversal_en    bool      `json:"Position_reversal_en"`    // 	"position_reversal_en": true,
	Delayed_stoploss_min    time.Time `json:"Delayed_stoploss_min"`    // 	"delayed_stoploss_min": "00:30:00",
	Stall_detect_period_min time.Time `json:"Stall_detect_period_min"` // 	"stall_detect_period_min": "00:30:00"
}

type Kite_Setting_S struct {
	Products     string `json:"Products"`
	Varieties    string `json:"Varieties"`
	OrderType    string `json:"OrderType"`
	Validities   string `json:"Validities"`
	PositionType string `json:"PositionType"`
}

type Trade_setting_S struct {
	TradeSimulate         bool    `json:"tradeSimulate"`
	OrderRoute            string  `json:"OrderRoute"`
	OptionLevel           int     `json:"OptionLevel"`
	OptionExpiryWeek      int     `json:"OptionExpiryWeek"`
	FuturesExpiryMonth    int     `json:"FuturesExpiryMonth"`
	SkipExipryWeekFutures bool    `json:"SkipExipryWeekFutures"`
	LimitAmount           float64 `json:"LimitAmount"`
}

type ControlData_S struct {
	Percentages     Percentage_S
	Target_Controls TargetControls_S
	Kite_Setting    Kite_Setting_S
	Trade_Setting   Trade_setting_S
}

type UserStrategies_S struct {
	Strategy     string
	Enabled      bool
	Engine       string
	Trigger_time time.Time
	Trigger_days string
	Cdl_size     int
	Instruments  string
	Controls     string
	CtrlData     ControlData_S
}

type Targets_S struct {
	Entry    float64 `json:"entry"`
	Target   float64 `json:"target"`
	Stoploss float64 `json:"stoploss"`
}

type Info_S struct {
	TradingSymbol     string  `json:"trading_symbol"`
	Exchange          string  `json:"exchange"`
	OrderIdEntr       uint64  `json:"order_id_entr"`
	OrderIdExit       uint64  `json:"order_id_exit"`
	QtyReq            float64 `json:"qty_req"`
	QtyFilled         float64 `json:"qty_filled"`
	UserExitRequested bool    `json:"user_exit_requested"`
	AvgPriceEnter     float64 `json:"avg_price_entr"`
	AvgPriceExit      float64 `json:"avg_price_exit"`
}

type OrderBook_S struct {
	Id            uint16
	Date          time.Time
	Instr         string
	Strategy      string
	Status        string
	Dir           string
	Exit_reason   string
	Info          Info_S
	Targets       Targets_S
	Orders_entr   []kiteconnect.Trade
	Orders_exit   []kiteconnect.Trade
	Post_analysis string
}

type ApiSignal struct {
	Status   string    `json:"stat us"`
	Id       uint16    `json:"id"`
	Date     time.Time `json:"date"`
	Instr    string    `json:"instr"`
	Strategy string    `json:"strategy"`
	Dir      string    `json:"dir"`
	Entry    float64   `json:"entry"`
	Target   float64   `json:"target"`
	Stoploss float64   `json:"stoploss"`
}

// Env variables required
var UserSettings = []string{
	"APP_LIVE_TRADING_MODE",
	"ZERODHA_USER_ID",
	"ZERODHA_PASSWORD",
	"ZERODHA_API_KEY",
	"ZERODHA_PIN",
	"ZERODHA_API_SECRET",
	"ZERODHA_TOTP_SECRET_KEY",
	"ZERODHA_REQ_TOKEN_URL",
	"TIMESCALEDB_ADDRESS",
	"TIMESCALEDB_USERNAME",
	"TIMESCALEDB_PASSWORD",
	"TIMESCALEDB_PORT",
	"ALGO_ANALYSIS_ADDRESS",
	"DB_TBL_TICK_NSEFUT",
	"DB_TBL_TICK_NSESTK",
	"DB_TBL_USER_SYMBOLS",
	"DB_TBL_USER_SETTING",
	"DB_TBL_USER_STRATEGIES",
	"DB_TBL_ORDER_BOOK",
	"DB_TEST_PREFIX",
	"DB_TBL_PREFIX_USER_ID",
}

var (
	Env = make(map[string]string)
)
