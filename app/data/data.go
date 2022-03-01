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
	Id             uint16    // 1
	Date           time.Time // 2
	Instr          string    // 3
	Strategy       string    // 4
	Dir            string    // 5
	Entry          float64   // 6
	Entry_time     time.Time // 7
	Target         float64   // 8
	Stoploss       float64   // 9
	Trade_id       uint64    // 10
	Exit_val       float64   // 11
	Exit_time      time.Time // 12
	Exit_reason    string    // 13
	Swing_min      float64   // 14
	Swing_max      float64   // 15
	Swing_min_time time.Time // 16
	Swing_max_time time.Time // 17
}
