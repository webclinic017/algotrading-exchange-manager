package trademgr

var startTrader_TblOdrbook_deleteAll = `DELETE FROM public.paragvb_order_book_test;`
var startTrader_TblUserStrategies_deleteAll = `DELETE FROM public.paragvb_strategies_test;`

var startTrader_TblUserStrategies_setup = `
INSERT INTO public.paragvb_strategies_test (strategy,enabled,engine,trigger_time,trigger_days,cdl_size,instruments,parameters) VALUES
	 ('S999-CONT-001',true,'IntraDay_DNP','%TRIGGERTIME','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'TT_TEST1','{
    "kite_setting": {
        "products": "MIS",
        "varieties": "regular",
        "order_type": "MARKET",
        "validities": "IOC",
        "position_type": "day"
    },
    "controls": {
        "trade_simulate": true,
        "target_per": 11,
        "stoploss_per": 21,
        "deep_stoploss_per": 31,
        "delayed_stoploss_min": "2018-09-22T23:23:23Z",
        "stall_detect_period_min": "2018-09-22T22:22:22Z",
        "budget_max_per": 51,
        "limit_amount": 30001,
        "trail_target_en": true,
        "position_reversal_en": true,
        "winning_ratio": 81
    },
    "options_setting": {
        "order_route": "option-buy",
        "option_level": -1,
        "option_expiry_week": 0
    },
    "futures_setting": {
        "futures_expiry_month": 0,
        "skip_exipry_week": true
    }
}'),
	 ('S999-TEST-002',true,'IntraDay_DNP','%TRIGGERTIME','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'TT_TEST2','{
    "kite_setting": {
        "products": "MIS",
        "varieties": "regular",
        "order_type": "MARKET",
        "validities": "IOC",
        "position_type": "day"
    },
    "controls": {
        "trade_simulate": true,
        "target_per": 12,
        "stoploss_per": 22,
        "deep_stoploss_per": 32,
        "delayed_stoploss_min": "2018-09-22T23:23:23Z",
        "stall_detect_period_min": "2018-09-22T22:22:22Z",
        "budget_max_per": 52,
        "limit_amount": 30002,
        "trail_target_en": false,
        "position_reversal_en": false,
        "winning_ratio": 82
    },
    "options_setting": {
        "order_route": "option-buy",
        "option_level": -1,
        "option_expiry_week": 0
    },
    "futures_setting": {
        "futures_expiry_month": 0,
        "skip_exipry_week": true
    }
}');
`

var startTrader_TblUserStrategies_EqRelianceREAL = `
INSERT INTO public.paragvb_strategies_test (strategy,enabled,engine,trigger_time,trigger_days,cdl_size,instruments,parameters) VALUES
	('S999-TEST-002',true,'IntraDay_DNP','%TRIGGERTIME','Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday',1,'RELIANCE','{
    "kite_setting": {
        "products": "MIS",
        "varieties": "regular",
        "order_type": "MARKET",
        "validities": "IOC",
        "position_type": "day"
    },
    "controls": {
        "trade_simulate": true,
        "target_per": 12,
        "stoploss_per": 22,
        "deep_stoploss_per": 32,
        "delayed_stoploss_min": "2018-09-22T23:23:23Z",
        "stall_detect_period_min": "2018-09-22T22:22:22Z",
        "budget_max_per": 52,
        "limit_amount": 30002,
        "trail_target_en": false,
        "position_reversal_en": false,
        "winning_ratio": 82
    },
    "options_setting": {
        "order_route": "option-buy",
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
('2022-04-20',  'CONTINOUS_test1',  'S99-TEST-002', 'AwaitSignal',  '',     '',             '{}','{}','[{}]','[{}]','{}'),
('2022-04-20',  'TIMETRIG_test2',   'S999-TEST-002', 'AwaitSignal',  '',     '',             '{}','{}','[{}]','[{}]','{}')
;`
