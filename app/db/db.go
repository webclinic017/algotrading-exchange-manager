package db

import (
	"context"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dbpool *pgxpool.Pool
var dbTick []kite.TickData

func DbInit() bool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	ctx := context.Background()
	var err error

	dbUrl := os.Getenv("DATABASE_URL")
	dbpool, err = pgxpool.Connect(ctx, dbUrl)
	myCon, _ := dbpool.Acquire(ctx)
	defer myCon.Release()

	if err != nil {
		srv.ErrorLogger.Printf("Unable to connect to database: %v\n", err)
		return false
	}

	var greeting string
	err = myCon.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)

	if err != nil {
		srv.ErrorLogger.Printf("QueryRow failed: %v\n", err)
		return false
	}
	srv.InfoLogger.Printf("connected to DB : " + greeting)

	_, table_check := myCon.Query(ctx, "select * from "+"zerodha_ticks"+";")

	if table_check != nil {
		srv.InfoLogger.Printf("DB Does not exist, creating now!: %v\n", err)

		// check if table exist, else create it
		queryCreateTicksTable := `CREATE TABLE 
								zerodha_ticks (
									time TIMESTAMP NOT NULL,
									symbol VARCHAR(30) NOT NULL,
									last_traded_price double precision NOT NULL,
									buy_demand bigint NOT NULL,
									sell_demand bigint NOT NULL,
									trades_till_now bigint NOT NULL,
									open_interest bigint NOT NULL
								);
						SELECT create_hypertable('zerodha_ticks', 'time');
						SELECT set_chunk_time_interval('zerodha_ticks', INTERVAL '24 hours');
						`

		//execute statement, fails if table already exists
		_, err = myCon.Exec(ctx, queryCreateTicksTable)
		if err != nil {
			srv.WarningLogger.Printf("DB CREATE: %v\n", err)
		}
		createViews()
		setupDbCompression()

	}
	// check if table exist, else create it
	_, table_check = myCon.Query(ctx, "select * from "+"zerodha_ticks_id_daily"+";")

	if table_check != nil {
		queryCreateSymbolsTable := `CREATE TABLE 
									zerodha_ticks_id_daily (
									time TIMESTAMP NULL,
									nse_symbol VARCHAR(30) NULL,
									mcx_symbol VARCHAR(30) NULL									
								);
						SELECT create_hypertable('zerodha_ticks_id_daily', 'time');
						SELECT set_chunk_time_interval('zerodha_ticks_id_daily', INTERVAL '1 YEAR');
						`

		//execute statement, fails if table already exists
		_, err = myCon.Exec(ctx, queryCreateSymbolsTable)
		if err != nil {
			srv.WarningLogger.Printf("DB CREATE: %v\n", err)
		}
	}

	return true
}

func setupDbCompression() {

	ctx := context.Background()
	myCon, _ := dbpool.Acquire(ctx)
	defer myCon.Release()

	_, err := myCon.Exec(ctx, `ALTER TABLE zerodha_ticks SET (
									timescaledb.compress,
									timescaledb.compress_segmentby = 'symbol'); 
								
									SELECT add_compression_policy('zerodha_ticks', INTERVAL '7 days');
								`)
	if err != nil {
		srv.WarningLogger.Printf("Error setting up DB Compression: %v\n", err)
	}
}

func createViews() {

	ctx := context.Background()
	myCon, _ := dbpool.Acquire(ctx)
	defer myCon.Release()

	_, err := myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_1min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('1 minutes', time) AS bucket, 
									symbol,
									FIRST(time, time) as first_time,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume,
									LAST(time, time) as last_time
								FROM
									zerodha_ticks
								GROUP by
									symbol, bucket
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_1min',
									start_offset => INTERVAL '1 days',
									end_offset => NULL,
									schedule_interval => INTERVAL '1 minutes');
								`)
	if err != nil {
		srv.WarningLogger.Printf("Error creating candles_1min: %v\n", err)
	}

	_, err = myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_3min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('3 minutes', time) AS bucket, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									zerodha_ticks
								GROUP by
									symbol, bucket
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_3min',
									start_offset => INTERVAL '1 days',
									end_offset => NULL,
									schedule_interval => INTERVAL '3 minutes');
	`)
	if err != nil {
		srv.WarningLogger.Printf("Error creating candles_3min: %v\n", err)
	}

	_, err = myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_5min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('5 minutes', time) AS bucket, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									zerodha_ticks
								GROUP by
									symbol, bucket
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_5min',
									start_offset => INTERVAL '1 days',
									end_offset => NULL,
									schedule_interval => INTERVAL '5 minutes');
	`)
	if err != nil {
		srv.WarningLogger.Printf("Error creating candles_5min: %v\n", err)
	}

	_, err = myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_10min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('10 minutes', time) AS bucket, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									zerodha_ticks
								GROUP by
									symbol, bucket
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_10min',
									start_offset => INTERVAL '1 days',
									end_offset => NULL,
									schedule_interval => INTERVAL '10 minutes');
	`)
	if err != nil {
		srv.WarningLogger.Printf("Error creating candles_10min: %v\n", err)
	}

	_, err = myCon.Exec(ctx, `CREATE MATERIALIZED VIEW candles_15min
								WITH (timescaledb.continuous) AS
								SELECT time_bucket('15 minutes', time) AS bucket, 
									symbol,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
								FROM
									zerodha_ticks
								GROUP by
									symbol, bucket
								WITH NO DATA;

								SELECT add_continuous_aggregate_policy('candles_15min',
									start_offset => INTERVAL '1 days',
									end_offset => NULL,
									schedule_interval => INTERVAL '15 minutes');
	`)
	if err != nil {
		srv.WarningLogger.Printf("Error creating candles_15min: %v\n", err)
	}
}

func StoreTickInDb() {
	for v := range kite.ChTick { // read from tick channel
		dbTick = append(dbTick, v)
		if len(dbTick) > 250 {
			go executeBatch(dbTick)
			dbTick = nil
		}
	}
	go executeBatch(dbTick)
	dbTick = nil
}

func StoreSymbolsInDb(nse_symbol string, mcx_symbol string) {
	ctx := context.Background()
	myCon, _ := dbpool.Acquire(ctx)
	defer myCon.Release()

	timestamp := time.Now()
	queryInsertMetadata := `INSERT INTO zerodha_ticks_id_daily (
		time,
		nse_symbol,
		mcx_symbol)
		VALUES
		($1, $2, $3);`

	_, err := myCon.Exec(ctx, queryInsertMetadata,
		timestamp,
		nse_symbol,
		mcx_symbol)
	if err != nil {
		srv.ErrorLogger.Printf("Unable to insert data into 'symbol ID' database: %v\n", err)
	}
}
func CloseDBpool() {

}

func executeBatch(dataTick []kite.TickData) {
	// start := time.Now()

	defer func() {
		if err := recover(); err != nil {
			srv.WarningLogger.Printf("DB Not intialised: ", err)
		}
	}()

	batch := &pgx.Batch{}

	queryInsertTimeseriesData := `INSERT INTO zerodha_ticks (
					time,
					symbol,
					last_traded_price,
					buy_demand, sell_demand,
					trades_till_now,
					open_interest)
					VALUES
					($1, $2, $3, $4, $5, $6, $7);`

	for i := range dataTick {
		var ct kite.TickData = dataTick[i]

		batch.Queue(queryInsertTimeseriesData,
			ct.Timestamp,
			ct.Symbol,
			ct.LastTradedPrice,
			ct.Buy_Demand,
			ct.Sell_Demand,
			ct.TradesTillNow,
			ct.OpenInterest)
	}

	ctx := context.Background()

	myCon, _ := dbpool.Acquire(ctx)
	defer myCon.Release()

	// stat := dbpool.Stat()

	br := myCon.SendBatch(ctx, batch)
	_, err := br.Exec()

	if err != nil {
		srv.WarningLogger.Printf("Unable to execute statement in batch queue %v\n", err)
	}

	// elapsed := time.Since(start)
	// fmt.Printf("AcquiredConns: %d TotalConns: %d DB operations took %s\n", stat.AcquiredConns(), stat.TotalConns(), elapsed)

}
