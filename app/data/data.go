package data

import (
	"time"
)

type Percentage struct {
	Target float64 // "target": 1,
	SL     float64 // "sl": 1,
	DeepSL float64 // "deepsl": 1
}

type TargetControls struct {
	Trail_target_en         bool      // 	"trail_target_en": true,
	Position_reversal_en    bool      // 	"position_reversal_en": true,
	Delayed_stoploss_min    time.Time // 	"delayed_stoploss_min": "00:30:00",
	Stall_detect_period_min time.Time // 	"stall_detect_period_min": "00:30:00"
}

type ControlParams struct {
	Percentages     Percentage
	Target_Controls TargetControls
	TradeBase       string // 	"trade_base": "future","stock","option"
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
