package db

import (
	"context"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

func DbInit() bool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	ctx := context.Background()
	var err error

	dbUrl := os.Getenv("DATABASE_URL")
	dbPool, err = pgxpool.Connect(ctx, dbUrl)

	if err != nil {
		srv.ErrorLogger.Printf("Unable to connect to database: %v\n", err)
		return false
	}

	var greeting string
	err = dbPool.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)

	if err != nil {
		srv.ErrorLogger.Printf("QueryRow failed: %v\n", err)
		return false
	}
	srv.InfoLogger.Printf("connected to DB : " + greeting)

	// check if table exist, else create it
	queryCreateTicksTable := `CREATE TABLE 
								zerodha_nse_mcx_ticks (
									time TIMESTAMP NOT NULL,
									symbol VARCHAR(30) NOT NULL,
									last_traded_price double precision NOT NULL,
									buy_demand int NOT NULL,
									sell_demand int NOT NULL,
									last_traded_quantity int NOT NULL,
									open_interest int NOT NULL
								);
						SELECT create_hypertable('zerodha_nse_mcx_ticks', 'time');
						SELECT set_chunk_time_interval('zerodha_nse_mcx_ticks', INTERVAL '24 hours');
						`

	//execute statement, fails if table already exists
	_, err = dbPool.Exec(ctx, queryCreateTicksTable)
	if err != nil {
		srv.WarningLogger.Printf("DB CREATE: %v\n", err)
	}

	// check if table exist, else create it
	queryCreateSymbolsTable := `CREATE TABLE 
									zerodha_nse_mcx_symbols_id_daily (
									time TIMESTAMP NULL,
									nse_symbol VARCHAR(30) NULL,
									mcx_symbol VARCHAR(30) NULL									
								);
						SELECT create_hypertable('zerodha_nse_mcx_symbols_id_daily', 'time');
						SELECT set_chunk_time_interval('zerodha_nse_mcx_symbols_id_daily', INTERVAL '1 YEAR');
						`

	//execute statement, fails if table already exists
	_, err = dbPool.Exec(ctx, queryCreateSymbolsTable)
	if err != nil {
		srv.WarningLogger.Printf("DB CREATE: %v\n", err)
	}
	return true
}

func OptimiseDbSettings() {

	ctx := context.Background()
	queryCreateAggregate := `CREATE MATERIALIZED VIEW one_min_candles_time
								WITH (timescaledb.continuous) AS
								select time_bucket('1 minutes', time) AS bucket, 
									symbol,
									FIRST(time, time) as first_time,
									FIRST(last_traded_price, time) as open,
									MAX(last_traded_price) as high,
									MIN(last_traded_price) as low,
									LAST(last_traded_price, time) as close,
									LAST(time, time) as last_time
								from
									zerodha_nse_mcx_ticks
								GROUP by
									symbol, bucket;
								`
	_, err := dbPool.Exec(ctx, queryCreateAggregate)
	if err != nil {
		srv.WarningLogger.Printf("Unable to create 1-Min Table (Conitnous Aggregate) DB: %v\n", err)
	}

	queryAddCaPolicy := `SELECT add_continuous_aggregate_policy('one_min_candles_time',
									start_offset => NULL,
									end_offset => NULL,
									schedule_interval => INTERVAL '1 minutes');`

	_, err = dbPool.Exec(ctx, queryAddCaPolicy)
	if err != nil {
		srv.WarningLogger.Printf("Unable to create 1-Min Table (Conitnous Aggregate) DB: %v\n", err)
	}
}

func StoreTickInDb() {

	/*
		Timestamp:          tick.Timestamp.Time,
		Symbol:             InsNamesMap[fmt.Sprint(tick.InstrumentToken)],
		LastTradedPrice:    tick.LastPrice,
		Buy_Demand:         tick.TotalBuyQuantity,
		Sell_Demand:        tick.TotalSellQuantity,
		LastTradedQuantity: tick.LastTradedQuantity,
		OpenInterest:       tick.OI
	*/

	for v := range kite.ChTick {
		// fmt.Println("\nkite ch data rx ", v)
		//fmt.Println("Timestamp: ", v.Timestamp)
		ctx := context.Background()
		// kite.ChTick <- kite.TickData{Timestamp: "2021-11-30 22:12:10", Insttoken: 1, Lastprice: 1, Open: 1.1, High: 1.2, Low: 1.3, Close: 1.4, Volume: 9}

		queryInsertMetadata := `INSERT INTO zerodha_nse_mcx_ticks (
			time,
			symbol, 
			last_traded_price,
			buy_demand, sell_demand, 
			last_traded_quantity,
			open_interest)
			VALUES 
			($1, $2, $3, $4, $5, $6, $7);`

		_, err := dbPool.Exec(ctx, queryInsertMetadata,
			v.Timestamp,
			v.Symbol,
			v.LastTradedPrice,
			v.Buy_Demand, v.Sell_Demand,
			v.LastTradedQuantity,
			v.OpenInterest)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to insert data into database: %v\n", err)
		}
	}
}

func StoreSymbolsInDb(nse_symbol string, mcx_symbol string) {
	ctx := context.Background()
	timestamp := time.Now()
	queryInsertMetadata := `INSERT INTO zerodha_nse_mcx_symbols_id_daily (
		time,
		nse_symbol, 
		mcx_symbol)
		VALUES 
		($1, $2, $3);`

	_, err := dbPool.Exec(ctx, queryInsertMetadata,
		timestamp,
		nse_symbol,
		mcx_symbol)
	if err != nil {
		srv.ErrorLogger.Printf("Unable to insert data into 'symbol ID' database: %v\n", err)
	}

}

func CloseDBPool() {
	dbPool.Close()
}
