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
