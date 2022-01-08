2021-10-24 12: 06: 08.991 Ticks: [
    {'tradable': True, 'mode': 'ltp', 'instrument_token': 18257666, 'last_price': 40670.0
    },
    {'tradable': True, 'mode': 'ltp', 'instrument_token': 18332418, 'last_price': 2651.8
    }
]

{Mode: "full", InstrumentToken: 18257666, IsTradable: true, IsIndex: false, Timestamp: github.com/zerodha/gokiteconnect/v4/models.Time {Time: (*time.Time)(0xc0000d99b8)
    }, LastTradeTime: github.com/zerodha/gokiteconnect/v4/models.Time {Time: (*time.Time)(0xc0000d99d0)
    }, LastPrice: 40670, LastTradedQuantity: 0, TotalBuyQuantity: 0, TotalSellQuantity: 0, VolumeTraded: 0, TotalBuy: 0, TotalSell: 0, AverageTradePrice: 0, OI: 0, OIDayHigh: 0, OIDayLow: 0, NetChange: 252.6999999999971, OHLC: github.com/zerodha/gokiteconnect/v4/models.OHLC {InstrumentToken: 0, Open: 40600, High: 40820, Low: 40465, Close: 40417.3
    }
}

Tick: {full 18257666 true false 2021-10-22 16: 42: 11 +0530 IST 2021-10-22 15: 29: 45 +0530 IST 40670 0 0 0 0 0 0 0 0 0 0 252.6999999999971 {
        0 40600 40820 40465 40417.3
    } {
        [
            {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            }
        ] [
            {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            }
        ]
    }
}

Tick: {full 18332418 true false 2021-10-22 16: 42: 12 +0530 IST 2021-10-22 15: 29: 57 +0530 IST 2651.8 0 0 0 0 0 0 0 0 0 0 -3.099999999999909 {
        0 2655.7 2695 2643.5 2654.9
    } {
        [
            {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            }
        ] [
            {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            } {
                0 0 0
            }
        ]
    }
}
Time: 2021-10-22 16: 42: 12 +0530 IST
Instrument: 18332418
LastPrice: 2651.8
Open: 2655.7
High: 2695
Low: 2643.5
Close: 2654.9
Volumne: 0


Connected
Tick:  {full 256265 false true 2021-11-30 17:31:31 +0530 IST 0001-01-01 00:00:00 +0000 UTC 16983.2 0 0 0 0 0 0 0 0 0 0 -70.75 {0 17051.15 17324.65 16931.4 17053.95} {[{0 0 0} {0 0 0} {0 0 0} {0 0 0} {0 0 0}] [{0 0 0} {0 0 0} {0 0 0} {0 0 0} {0 0 0}]}}
Tick:  {full 18332418 true false 2021-11-30 16:35:21 +0530 IST 2021-11-30 15:29:58 +0530 IST 2414 250 150250 1916000 12010000 0 0 2434.55 35591250 35591250 33723750 -31.59999999999991 {0 2463.95 2479 2401 2445.6} {[{0 0 1} {2411.35 750 1} {2411.1 250 1} {2411 1250 4} {2410.6 1000 2}] [{0 0 1} {2414.8 250 1} {2415 750 2} {2415.75 250 1} {2415.8 250 1}]}}
Tick:  {full 18257666 true false 2021-11-30 16:35:20 +0530 IST 2021-11-30 15:29:59 +0530 IST 35711 25 52000 147200 5527825 0 0 36249.32 2682150 2757875 2512550 -389.4499999999971 {0 36144.75 36849 35645.7 36100.45} {[{0 0 1} {35695 50 1} {35694.5 25 1} {35694 25 1} {35692.35 250 1}] [{0 0 1} {35715.5 25 1} {35716 125 1} {35716.1 25 1} {35717.25 250 1}]}}

kite ch data rx  {2021-11-30 17:31:31 +0530 IST 16983.2 256265 17051.15 17324.65 16931.4 17053.95 0}

kite ch data rx  {2021-11-30 16:35:21 +0530 IST 2414 18332418 2463.95 2479 2401 2445.6 12010000}

kite ch data rx  {2021-11-30 16:35:20 +0530 IST 35711 18257666 36144.75 36849 35645.7 36100.45 5527825}

                                                                                                                                                                                                                                                                        	// fmt.Printf("AcquireCount: %d\n", stat.AcquireCount())
	// fmt.Printf("AcquireDuration: %d\n", stat.AcquireDuration())
	fmt.Printf("AcquiredConns: %d\n", stat.AcquiredConns())
	// fmt.Printf("CanceledAcquireCount: %d\n", stat.CanceledAcquireCount())
	// fmt.Printf("ConstructingConns: %d\n", stat.ConstructingConns())
	// fmt.Printf("EmptyAcquireCount: %d\n", stat.EmptyAcquireCount())
	// fmt.Printf("IdleConns: %d\n", stat.IdleConns())
	// fmt.Printf("MaxConns: %d\n", stat.MaxConns())
	fmt.Printf("TotalConns: %d\n", stat.TotalConns())                                                                                                                                                       

{
	{ 
		true 45454.7 
		{
			0 
			45454.7 
			0 
			0 
			45454.7 
			45454.7
		}
		{
			0
			0 
			0 
			0 
			0 
			0 
			0 
			0 
			0 
			0 
			0 
			0
		}
	}
	{
		true
		0
		{
			0
			0
			0
			0
			0
			0
		}
		{
			0
			0
			0
			0
			0
			0
			0
			0
			0
			0
			0
			0
		}
	}
}                                                                                         







