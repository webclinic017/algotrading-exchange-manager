package appdata

import (
	"time"
)

// Global variables
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

type Percentage struct {
	Target      float64 // "target": 1,
	SL          float64 // "sl": 1,
	DeepSL      float64 // "deepsl": 1
	MaxBudget   float64 // "limit_budget": 50%,
	WinningRate float64 // "winning_rate": 80%,
}

type TargetControls struct {
	Trail_target_en         bool      // 	"trail_target_en": true,
	Position_reversal_en    bool      // 	"position_reversal_en": true,
	Delayed_stoploss_min    time.Time // 	"delayed_stoploss_min": "00:30:00",
	Stall_detect_period_min time.Time // 	"stall_detect_period_min": "00:30:00"
}

type Kite_Setting struct {
	Products     string
	Varieties    string
	OrderType    string
	Validities   string
	PositionType string
}

type Trade_setting struct {
	OrderRoute            string
	OptionLevel           int
	OptionExpiryWeek      int
	FuturesExpiryMonth    int
	SkipExipryWeekFutures bool
	LimitAmount           float64
}

type ControlParams struct {
	Percentages     Percentage
	Target_Controls TargetControls
	KiteSettings    Kite_Setting
	TradeSettings   Trade_setting
}

type Strategies struct {
	Strategy     string    // 0
	Enabled      bool      // 1
	Engine       string    // 2
	Trigger_time time.Time // 3
	Trigger_days string    // 4
	Cdl_size     int       // 6
	Instruments  string    // 7
	Controls     string
	CtrlParam    ControlParams
}

type ApiSignal struct {
	Status   string    `json:"status"`
	Id       uint16    `json:"id"`
	Date     time.Time `json:"date"`
	Instr    string    `json:"instr"`
	Strategy string    `json:"strategy"`
	Dir      string    `json:"dir"`
	Entry    float64   `json:"entry"`
	Target   float64   `json:"target"`
	Stoploss float64   `json:"stoploss"`
}

type TradeSignal struct {
	Id                 uint16    // 1
	Date               time.Time // 2
	Instr              string    // 3
	Strategy           string    // 4
	Status             string    // 5
	Instr_id           int       // 6
	Dir                string    // 6
	Entry              float64   //
	Target             float64   //
	Stoploss           float64   //
	Order_id           uint64    //
	Order_trades_entry string
	Order_trades_exit  string
	Order_simulation   string
	Exit_reason        string
	Post_analysis      string
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
