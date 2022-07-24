package db

import (
	"algo-ex-mgr/app/appdata"
	"strings"
)

//  ---------------------------------- CREATE TABLES  ----------------------------------

var DB_EXISTS_QUERY = "SELECT datname FROM pg_catalog.pg_database  WHERE lower(datname) = lower('algotrading');"
var DB_CREATE_QUERY = "CREATE DATABASE algotrading;"
var DB_TRADEMGR_EXISTS_QUERY = "SELECT controls from %DB_TBL_USER_SETTING WHERE name = 'trademgr.exits';"

var DB_CREATE_TBL_INSTRUMENTS = `DROP TABLE IF EXISTS %DB_TBL_INSTRUMENTS;

								CREATE TABLE %DB_TBL_INSTRUMENTS (
										instrument_token int8 NULL,
										exchange_token int8 NULL,
										tradingsymbol text NULL,
										"name" text NULL,
										last_price int8 NULL,
										expiry text NULL,
										strike float8 NULL,
										tick_size float8 NULL,
										lot_size int8 NULL,
										instrument_type text NULL,
										segment text NULL,
										exchange text NULL
									);`

var DB_CREATE_TABLE_TICKER_NSEFUT = `CREATE TABLE %DB_TBL_TICK_NSEFUT
							 		(
										time TIMESTAMP NOT NULL,
										symbol VARCHAR(30) NOT NULL,
										last_traded_price double precision NOT NULL DEFAULT 0,
										buy_demand bigint NOT NULL DEFAULT 0,
										sell_demand bigint NOT NULL DEFAULT 0,
										trades_till_now bigint NOT NULL DEFAULT 0,
										open_interest bigint NOT NULL DEFAULT 0
									);
								SELECT create_hypertable('%DB_TBL_TICK_NSEFUT', 'time');
								SELECT set_chunk_time_interval('%DB_TBL_TICK_NSEFUT', INTERVAL '7 days');`

var DB_CREATE_TABLE_TICKER_NSESTK = `CREATE TABLE %DB_TBL_TICK_NSESTK
									(
									   time TIMESTAMP NOT NULL,
									   symbol VARCHAR(30) NOT NULL,
									   last_traded_price double precision NOT NULL DEFAULT 0,
									   buy_demand bigint NOT NULL DEFAULT 0,
									   sell_demand bigint NOT NULL DEFAULT 0,
									   trades_till_now bigint NOT NULL DEFAULT 0,
									   open_interest bigint NOT NULL DEFAULT 0
								   );
							   SELECT create_hypertable('%DB_TBL_TICK_NSESTK', 'time');
							   SELECT set_chunk_time_interval('%DB_TBL_TICK_NSESTK', INTERVAL '7 days');`

var DB_CREATE_TABLE_USER_SYMBOLSwDel = `DROP TABLE IF EXISTS %DB_TBL_USER_SYMBOLS;
										` + DB_CREATE_TABLE_USER_SYMBOLS

var DB_CREATE_TABLE_USER_SYMBOLS = `CREATE TABLE %DB_TBL_USER_SYMBOLS (
									symbol varchar NOT NULL,
									track bool NULL DEFAULT false,
									segment varchar NOT NULL,
									mysymbol varchar NOT NULL,
									strikestep float4 NULL DEFAULT 0,
									exchange varchar NULL
								);`

var DB_CREATE_TABLE_USER_SETTING = `CREATE TABLE %DB_TBL_USER_SETTING (
									name varchar NOT NULL,
									controls JSON NOT NULL
								);`

var DB_CREATE_TABLE_USER_STRATEGIES = `CREATE TABLE %DB_TBL_USER_STRATEGIES (
										strategy VARCHAR(100) UNIQUE NOT NULL,
										enabled BOOLEAN NOT NULL DEFAULT 'false',
										engine  VARCHAR(50) NOT NULL,
										trigger_time TIME NOT NULL,
										trigger_days VARCHAR(100) NOT NULL,
										cdl_size SMALLINT NOT NULL,
										instruments TEXT,
										controls JSON
									);`

var DB_CREATE_TABLE_ORDER_BOOK = `CREATE TABLE %DB_TBL_ORDER_BOOK (
									id SERIAL PRIMARY KEY NOT NULL,
									date DATE NOT NULL,
									instr TEXT NOT NULL,
									strategy  VARCHAR(100) NOT NULL,
									status TEXT,
									dir VARCHAR(50),
									exit_reason TEXT  DEFAULT 'NA',
									info JSON,
									targets JSON,
									api_signal_entr JSON,
									api_signal_exit JSON,
									orders_entr JSON,
									orders_exit JSON,
									post_analysis JSON
								);`

//  ---------------------------------- COMPRESSION ----------------------------------

var DB_NSEFUT_COMPRESSION_QUERY = `ALTER TABLE $1 SET 
							(
								timescaledb.compress,
								timescaledb.compress_segmentby = 'symbol'
							); 
							SELECT add_compression_policy('$1 ', INTERVAL '30 days');
						`

//  ---------------------------------- VIEWS ----------------------------------

var DB_VIEW_EXISTS = `
					SELECT view_name 
					FROM timescaledb_information.continuous_aggregates
					WHERE view_name = $1;`

var DB_VIEW_CREATE_FUT = `
					CREATE MATERIALIZED VIEW %DB_TBL_CDL_VIEW_FUT
					WITH (timescaledb.continuous) AS
					SELECT time_bucket('1 minutes', time) AS candle, 
						symbol,
						FIRST(last_traded_price, time) as open,
						MAX(last_traded_price) as high,
						MIN(last_traded_price) as low,
						LAST(last_traded_price, time) as close,
						ROUND(AVG(buy_demand),2) as buy_demand,
						ROUND(AVG(sell_demand),2) as sell_demand,
						ROUND(AVG(open_interest),2) as open_interest,
						LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
					FROM
						%DB_TBL_TICK_NSEFUT
					GROUP by
						symbol, candle
					WITH NO DATA;

					SELECT add_continuous_aggregate_policy('%DB_TBL_CDL_VIEW_FUT',
						start_offset => INTERVAL '3 day',
						end_offset => INTERVAL '1 day',
						schedule_interval => INTERVAL '1 day');
					`
var DB_VIEW_CREATE_STK = `
					CREATE MATERIALIZED VIEW %DB_TBL_CDL_VIEW_STK
					WITH (timescaledb.continuous) AS
					SELECT time_bucket('1 minutes', time) AS candle, 
						symbol,
						FIRST(last_traded_price, time) as open,
						MAX(last_traded_price) as high,
						MIN(last_traded_price) as low,
						LAST(last_traded_price, time) as close,
						ROUND(AVG(buy_demand),2) as buy_demand,
						ROUND(AVG(sell_demand),2) as sell_demand,
						ROUND(AVG(open_interest),2) as open_interest,
						LAST(trades_till_now, time) - FIRST(trades_till_now, time) as volume
					FROM
						%DB_TBL_TICK_NSESTK
					GROUP by
						symbol, candle
					WITH NO DATA;

					SELECT add_continuous_aggregate_policy('%DB_TBL_CDL_VIEW_STK',
						start_offset => INTERVAL '3 day',
						end_offset => INTERVAL '1 day',
						schedule_interval => INTERVAL '1 day');
					`

var sqlQueryViewStkGetID = `SELECT job_id 
							FROM timescaledb_information.jobs j, timescaledb_information.continuous_aggregates ca
							WHERE 
								j.hypertable_name = ca.materialization_hypertable_name
								AND
								ca.view_name  = '%DB_TBL_CDL_VIEW_STK';`

var sqlQueryViewFutGetID = `SELECT job_id 
							FROM timescaledb_information.jobs j, timescaledb_information.continuous_aggregates ca
							WHERE 
								j.hypertable_name = ca.materialization_hypertable_name
								AND
								ca.view_name  = '%DB_TBL_CDL_VIEW_FUT';`

// ---------------------------------- db-instruments ----------------------------------

var sqlQueryFutures = `SELECT i.instrument_token, ts.mysymbol
						FROM %DB_TBL_USER_SYMBOLS ts, %DB_TBL_INSTRUMENTS i
						WHERE 
								ts.symbol = i.name
							and 
								ts.segment = i.instrument_type 
							and 
								ts.exchange = i.exchange
							and 
								EXTRACT(MONTH FROM TO_DATE(i.expiry,'YYYY-MM-DD')) = EXTRACT(MONTH FROM current_date)+1;`

// RULE: FUT COntracts for next month, query return null when somedays are left in current month but contract expires. eg. 29Apr2022

var sqlInstrDataQueryOptn = `SELECT tradingsymbol, lot_size
								FROM %DB_TBL_USER_SYMBOLS ts, %DB_TBL_INSTRUMENTS i
								WHERE 
										i.exchange = 'NFO'
									and
										ts.symbol = i.name 
									and 
										mysymbol= $1 
									and
										strike >= ($2 + ($3*ts.strikestep) )
									and
										strike < ($2 + ts.strikestep + ($3*ts.strikestep) )
									and
										instrument_type = $4
									and
										expiry > $5
									and
										expiry < $6				
								ORDER BY 
									expiry asc
								LIMIT 10;`

var sqlInstrDataQueryEQ = `SELECT tradingsymbol, lot_size
								FROM %DB_TBL_USER_SYMBOLS ts, %DB_TBL_INSTRUMENTS i
								WHERE 
									ts.symbol = i.tradingsymbol 
								and
									ts.exchange = i.exchange 
								and 
									ts.segment = 'EQ'
								and 
									ts.symbol = $1 
								LIMIT 1;
								`

var sqlInstrDataQueryFUT = `SELECT tradingsymbol, lot_size
							FROM %DB_TBL_USER_SYMBOLS ts, %DB_TBL_INSTRUMENTS i
							WHERE 
									ts.symbol = i.name 
								and 
									mysymbol= $1
								and 
									expiry > $2
								and 
									expiry < $3
								and 
									instrument_type = 'FUT'
							LIMIT 10;`

var sqlQueryNseEqTokens = `SELECT i.instrument_token, ts.mysymbol
							FROM %DB_TBL_USER_SYMBOLS ts, %DB_TBL_INSTRUMENTS i
							WHERE 
									ts.symbol = i.tradingsymbol
								and 
									i.instrument_type = 'EQ'
								and 
									ts.exchange = i.exchange;`

// ---------------------------------- db-orderbook ----------------------------------

var sqlqueryOrderBookId = "SELECT * FROM %DB_TBL_ORDER_BOOK WHERE id = %d"
var sqlqueryAllOrderBookCondition = "SELECT * FROM %DB_TBL_ORDER_BOOK WHERE status %s '%s'"

var sqlCreateOrder = `INSERT INTO  %DB_TBL_ORDER_BOOK (
	date,
	instr,
	strategy,
	status,
	dir,
	exit_reason,
	info,
	targets,
	api_signal_entr,
	api_signal_exit,
	orders_entr,
	orders_exit,
	post_analysis)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`

var sqlUpdateOrder = ` UPDATE %DB_TBL_ORDER_BOOK SET
	date = $1,
	instr = $2,
	strategy = $3,
	status = $4,
	dir = $5,
	exit_reason = $6,
	info = $7,
	targets = $8,
	api_signal_entr = $9,
	api_signal_exit = $10,
	orders_entr = $11,
	orders_exit = $12,
	post_analysis = $13
	WHERE id = $14
	;`

var sqlOrderCount = `SELECT COUNT(*) FROM %DB_TBL_ORDER_BOOK
						WHERE  (
							instr = $1
						AND
							date = $2
						AND
							strategy = $3)`

var sqlOrderId = `
	   SELECT id FROM %DB_TBL_ORDER_BOOK
	   WHERE  (
	   		instr = $1
	   	AND
	   		date = $2
	   	AND
	   		strategy = $3)`

// ---------------------------------- query-resolver  ----------------------------------

func dbSqlQuery(query string) string {

	for key, val := range appdata.Env {
		query = strings.Replace(query, "%"+key, val, -1)
	}

	return query
}

// ---------------------------------- db.go ----------------------------------
var sqlSaveInstruments = `INSERT INTO %DB_TBL_INSTRUMENTS (
	instrument_token,
	exchange_token,
	tradingsymbol,
	"name",
	last_price,
	expiry,
	strike,
	tick_size,
	lot_size,
	instrument_type,
	segment,
	exchange)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

var sqlSaveUserSymbols = `INSERT INTO %DB_TBL_USER_SYMBOLS
	(
		symbol,
		track,
		segment,
		mysymbol,
		strikestep,
		exchange
	)
	VALUES
	(
		$1, $2, $3, $4, $5, $6
	);`
