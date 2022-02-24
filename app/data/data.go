package data

import "time"

type Strategies struct {
	Strategy_id               string
	Strategy_en               bool
	P_engine                  string
	P_trigger_time            time.Time
	P_trigger_days            string
	P_target_per              int
	P_candle_size             int
	P_stoploss_per            int
	P_deep_stoploss_per       int
	P_delayed_stoploss_min    time.Time
	P_stall_detect_period_min time.Time
	P_trail_target_en         bool
	P_position_reversal_en    bool
	P_trade_symbols           string
}

type TradeSignal struct {
	Strategy_id          string
	S_date               string
	S_direction          string
	T_entry              float64
	T_entry_time         string
	S_target             float64
	S_stoploss           float64
	T_trade_confirmed_en bool
	S_instr_token        string
	R_exit_val           float64
	R_exit_time          string
	R_exit_reason        string
	R_swing_min          float64
	R_swing_max          float64
	R_swing_min_time     string
	R_swing_max_time     string
}
