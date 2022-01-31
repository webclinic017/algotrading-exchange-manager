package db

import (
	"context"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"os"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool
var dbTick []kite.TickData

var DB_EXISTS_QUERY = "SELECT datname FROM pg_catalog.pg_database  WHERE lower(datname) = lower('algotrading');"
var DB_CREATE_QUERY = "CREATE DATABASE algotrading;"
var DB_TABLE_ID_DECODED_NAME = `token_id_decoded`
var DB_CREATE_TABLE_ID_DECODED = `CREATE TABLE token_id_decoded
								(
									time TIMESTAMP NOT NULL,
									nse_symbol VARCHAR(30),
									mcx_symbol VARCHAR(30)
								);
						`
var DB_TABLE_TICKER_NAME = `ticks_data`
var DB_CREATE_TABLE_TICKER = `CREATE TABLE 
								ticks_data
							 (
									time TIMESTAMP NOT NULL,
									symbol VARCHAR(30) NOT NULL,
									last_traded_price double precision NOT NULL DEFAULT 0,
									buy_demand bigint NOT NULL DEFAULT 0,
									sell_demand bigint NOT NULL DEFAULT 0,
									trades_till_now bigint NOT NULL DEFAULT 0,
									open_interest bigint NOT NULL DEFAULT 0
								);
							SELECT create_hypertable('ticks_data', 'time');
							SELECT set_chunk_time_interval('ticks_data', INTERVAL '7 days');
						`

var DB_COMPRESSION_QUERY = `ALTER TABLE ticks_data SET 
							(
								timescaledb.compress,
								timescaledb.compress_segmentby = 'symbol'
							); 
							SELECT add_compression_policy('ticks_data ', INTERVAL '30 days');
						`

func connectDB() bool {
	ctx := context.Background()
	dbUrl := "postgres://" + os.Getenv("TIMESCALEDB_USERNAME") + ":" + os.Getenv("TIMESCALEDB_PASSWORD") + "@" + os.Getenv("TIMESCALEDB_ADDRESS") + ":" + os.Getenv("TIMESCALEDB_PORT") + "/postgres"

	// Check if you can connect to DB server (accessing 'postgres' defualt DB)
	dbPoolDefault, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		srv.ErrorLogger.Println("Could not connect with 'postgres' DB\n", err)
		return false
	}
	defer dbPoolDefault.Close()

	myCon, err := dbPoolDefault.Acquire(ctx)
	defer myCon.Release()
	if err != nil {
		srv.ErrorLogger.Printf("Could not acquire Context, too many operations?: %v\n", err)
		return false
	}

	// check if 'algotrading' database exists, if not CREATE it
	var retVal string
	myCon.QueryRow(ctx, DB_EXISTS_QUERY).Scan(&retVal)

	if len(retVal) == 0 {
		srv.InfoLogger.Printf("algotrading DB Does not exist, creating now!: %v\n", err)

		//execute statement, fails if table already exists
		myCon2, _ := dbPoolDefault.Acquire(ctx)
		defer myCon.Release()
		_, err = myCon2.Exec(ctx, DB_CREATE_QUERY)
		if err != nil {
			srv.ErrorLogger.Printf("Failed to CREATE algotrading DB: %v\n", err)
			return false
		}
	}
	return true
}

func DbInit() bool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	srv.InfoLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Db Checks",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	ctx := context.Background()
	dbUrl := "postgres://" + os.Getenv("TIMESCALEDB_USERNAME") + ":" + os.Getenv("TIMESCALEDB_PASSWORD") + "@" + os.Getenv("TIMESCALEDB_ADDRESS") + ":" + os.Getenv("TIMESCALEDB_PORT") + "/algotrading"

	if connectDB() {
		// 1. Connect with 'algotrading' DB
		var err error
		dbPool, err = pgxpool.Connect(ctx, dbUrl)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to connect with 'algotrading db' %v\n", err)
			return false
		}
		// 2. Aquire context
		myCon, err := dbPool.Acquire(ctx)
		defer myCon.Release()
		if err != nil {
			srv.ErrorLogger.Printf("Could not acquire Context, too many operations?: %v\n", err)
			return false
		}

		// 3. Check if 'ticker' table exists, if not CREATE it
		if createTable(DB_TABLE_TICKER_NAME, DB_CREATE_TABLE_TICKER) {
			createViews()
			setupDbCompression()
			createTable(DB_TABLE_ID_DECODED_NAME, DB_CREATE_TABLE_ID_DECODED)
			srv.InfoLogger.Printf("DB checks completed\n")
			return true
		} else {
			return false
		}

	} else {
		return false
	}

}

func createTable(tblName string, sqlquery string) bool {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)

	var retVal string

	query := "select table_name from information_schema.tables WHERE table_name = '" + tblName + "';"
	myCon.QueryRow(ctx, query).Scan(&retVal)
	// if err != nil {
	// 	srv.WarningLogger.Printf("Failed to CREATE %s table : %v\n", tblName, err)
	// }

	if len(retVal) == 0 {
		srv.InfoLogger.Printf("%s Does not exist, creating now!\n", tblName)
		_, err := myCon.Exec(ctx, sqlquery)
		if err != nil {
			srv.WarningLogger.Printf("Failed to CREATE %s table : %v\n", tblName, err)
			myCon.Release()
			return false
		}
	}
	myCon.Release()
	return true
}

func setupDbCompression() {

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	_, _ = myCon.Exec(ctx, DB_COMPRESSION_QUERY)
	// if err != nil {
	// 	srv.WarningLogger.Printf("Error setting up DB Compression: %v\n", err)
	// }
}

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
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	timestamp := time.Now()
	queryInsertMetadata := `INSERT INTO token_id_decoded (
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

func CloseDBpool() bool {
	dbPool.Close()
	return false
}

func executeBatch(dataTick []kite.TickData) {

	defer func() {
		if err := recover(); err != nil {
			srv.WarningLogger.Print("DB Not intialised: ", err)
		}
	}()

	batch := &pgx.Batch{}

	queryInsertTimeseriesData := `INSERT INTO ticks_data
 (
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

	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	// stat := dbpool.Stat()

	br := myCon.SendBatch(ctx, batch)
	_, err := br.Exec()

	if err != nil {
		srv.WarningLogger.Printf("Unable to execute statement in batch queue %v\n", err)
	}
}
