package data

import "time"

type Strategies struct {
	Strategy_id               string    // 0
	Strategy_en               bool      // 1
	P_engine                  string    // 2
	P_trigger_time            time.Time // 3
	P_trigger_days            string    // 4
	P_target_per              int       // 5
	P_candle_size             int       // 6
	P_stoploss_per            int       // 7
	P_deep_stoploss_per       int       // 8
	P_delayed_stoploss_min    time.Time // 9
	P_stall_detect_period_min time.Time // 10
	P_trail_target_en         bool      // 11
	P_position_reversal_en    bool      // 12
	P_trade_symbols           string    // 13
}

type TradeSignal struct {
	S_order_id           uint16    // 0
	Strategy_id          string    // 1
	S_date               time.Time // 2
	S_direction          string    // 3
	T_entry              float64   // 4
	T_entry_time         time.Time // 5
	S_target             float64   // 6
	S_stoploss           float64   // 7
	T_trade_confirmed_en bool      // 8
	S_instr_token        string    // 9
	R_exit_val           float64   // 10
	R_exit_time          time.Time // 11
	R_exit_reason        string    // 12
	R_swing_min          float64   // 13
	R_swing_max          float64   // 14
	R_swing_min_time     time.Time // 15
	R_swing_max_time     time.Time // 16
}
