package db

import (
	"context"
	"goTicker/app/srv"

	"github.com/jackc/pgconn"
)

func createViews() {

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	_, err := myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_1min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('1 minutes', time) AS candle, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									ticks_data
								
								GROUP by
									symbol, candle
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_1min',
									start_offset => NULL,
									end_offset => INTERVAL '1 minutes',
									schedule_interval => INTERVAL '1 minutes');
								`)
	if err != nil {
		pgerr, _ := err.(*pgconn.PgError)
		if pgerr.Code != "42P07" {
			srv.WarningLogger.Printf("Error creating candles_1min: %v\n", err)
		}
	}

	_, err = myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_3min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('3 minutes', time) AS candle, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									ticks_data
								
								GROUP by
									symbol, candle
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_3min',
									start_offset => NULL,
									end_offset => INTERVAL '3 minutes',
									schedule_interval => INTERVAL '3 minutes');
	`)

	if err != nil {
		pgerr, _ := err.(*pgconn.PgError)
		if pgerr.Code != "42P07" {
			srv.WarningLogger.Printf("Error creating candles_3min: %v\n", err)
		}
	}

	_, err = myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_5min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('5 minutes', time) AS candle, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									ticks_data
								
								GROUP by
									symbol, candle
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_5min',
									start_offset => NULL,
									end_offset => INTERVAL '5 minutes',
									schedule_interval => INTERVAL '5 minutes');
	`)
	if err != nil {
		pgerr, _ := err.(*pgconn.PgError)
		if pgerr.Code != "42P07" {
			srv.WarningLogger.Printf("Error creating candles_5min: %v\n", err)
		}
	}
}
