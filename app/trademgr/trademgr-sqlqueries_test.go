package trademgr

var settings_exits_deleteAll = `UPDATE paragvb_setting_test set controls = '' WHERE name = 'trademgr.exits';`
var settings_exits_setVal = `UPDATE paragvb_setting_test set controls = '%EXIT_ID' WHERE name = 'trademgr.exits';`

var startTrader_TblOdrbook_deleteAll = `DELETE FROM public.paragvb_order_book_test;`
var startTrader_TblUserStrategies_deleteAll = `DELETE FROM public.paragvb_strategies_test;`

var startTrader_TblUserStrategies_setup = `
INSERT INTO public.paragvb_strategies_test (strategy,enabled,engine,trigger_time,trigger_days,cdl_size,instruments,parameters) VALUES
	 ('%STRATEGY_NAME_1',true,'IntraDay_DNP','%TRIGGERTIME','%TRIGGER_DAYS',1,'%SYMBOL_NAME_1','{
    "kite_setting": {
        "products": "MIS",
        "varieties": "regular",
        "order_type": "MARKET",
        "validities": "IOC",
        "position_type": "day",
        "order_route": "option-buy"
    },
    "controls": {
        "trade_simulate": true,
        "target_per": 11,
        "stoploss_per": 21,
        "deep_stoploss_per": 31,
        "delayed_stoploss_seconds":60,
        "stall_detect_period_seconds":60,
        "budget_max_per": 51,
        "limit_amount": 30001,
        "target_trail_enabled": true,
        "stoploss_trail_enabled": true,
        "position_reversal_en": true,
        "winning_ratio": 81
    },
    "options_setting": {
        "option_level": -1,
        "option_expiry_week": 0
    },
    "futures_setting": {
        "futures_expiry_month": 0,
        "skip_exipry_week": true
    }
}'),
	 ('%STRATEGY_NAME_2',true,'IntraDay_DNP','%TRIGGERTIME','%TRIGGER_DAYS',1,'%SYMBOL_NAME_2','{
    "kite_setting": {
        "products": "MIS",
        "varieties": "regular",
        "order_type": "MARKET",
        "validities": "IOC",
        "position_type": "day",
        "order_route": "option-buy"
    },
    "controls": {
        "trade_simulate": true,
        "target_per": 12,
        "stoploss_per": 22,
        "deep_stoploss_per": 32,
        "delayed_stoploss_seconds":60,
        "stall_detect_period_seconds":60,
        "budget_max_per": 52,
        "limit_amount": 30002,
        "target_trail_enabled": true,
        "stoploss_trail_enabled": true,
        "position_reversal_en": true,
        "winning_ratio": 82
    },
    "options_setting": {
        "option_level": -1,
        "option_expiry_week": 0
    },
    "futures_setting": {
        "futures_expiry_month": 0,
        "skip_exipry_week": true
    }
}');
`

var startTrader_TblUserStrategies_EqASHOKLEY_REAL = `
INSERT INTO public.paragvb_strategies_test (strategy,enabled,engine,trigger_time,trigger_days,cdl_size,instruments,parameters) VALUES
	('S999-TEST-002',true,'IntraDay_DNP','%TRIGGERTIME','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'ASHOKLEY','{
    "kite_setting": {
        "products": "MIS",
        "varieties": "regular",
        "order_type": "MARKET",
        "validities": "DAY",
        "position_type": "day",
        "order_route": "equity"
    },
    "controls": {
        "trade_simulate": false,
        "target_per": 12,
        "stoploss_per": 22,
        "deep_stoploss_per": 32,
        "delayed_stoploss_seconds":60,
        "stall_detect_period_seconds":60,
        "budget_max_per": 52,
        "limit_amount": 50,
        "target_trail_enabled": true,
        "stoploss_trail_enabled": true,
        "position_reversal_en": false,
        "winning_ratio": 82
    },
    "options_setting": {
        "option_level": -1,
        "option_expiry_week": 0
    },
    "futures_setting": {
        "futures_expiry_month": 0,
        "skip_exipry_week": true
    }
}');
`

var Test3_orderbook = `INSERT INTO public.paragvb_order_book_test 
("date",        instr,              strategy,       status,             dir,    exit_reason,    info,targets, orders_entr,orders_exit,post_analysis) VALUES
('2022-04-20',  'CONTINOUS_test1',  'S999-CONT-001', 'ExitOrdersPending',  'Bullish',     '',             '{"order_simulation":true}','{}','[{}]','[{}]','{}'),
('2022-04-20',  'CONTINOUS_test2',   'S999-CONT-001', 'ExitOrdersPending',  'Bullish',     '',             '{"order_simulation":true}','{}','[{}]','[{}]','{}')
;`
