package db

//  ---------------------------------- CREATE TABLES  ----------------------------------

var DB_EXISTS_QUERY = "SELECT datname FROM pg_catalog.pg_database  WHERE lower(datname) = lower('algotrading');"
var DB_CREATE_QUERY = "CREATE DATABASE algotrading;"
var DB_CREATE_TABLE_ID_DECODED = `CREATE TABLE token_id_decoded
								(
									time TIMESTAMP NOT NULL,
									nse_symbol VARCHAR(30),
									mcx_symbol VARCHAR(30)
								);
						`

var DB_CREATE_TABLE_TICKER = `CREATE TABLE $1
							 		(
										time TIMESTAMP NOT NULL,
										symbol VARCHAR(30) NOT NULL,
										last_traded_price double precision NOT NULL DEFAULT 0,
										buy_demand bigint NOT NULL DEFAULT 0,
										sell_demand bigint NOT NULL DEFAULT 0,
										trades_till_now bigint NOT NULL DEFAULT 0,
										open_interest bigint NOT NULL DEFAULT 0
									);
								SELECT create_hypertable('$1', 'time');
								SELECT set_chunk_time_interval('$1', INTERVAL '7 days');`

var DB_CREATE_TABLE_USER_SYMBOLS = `CREATE TABLE $1 (
									symbol varchar NOT NULL,
									track bool NULL DEFAULT false,
									segment varchar NOT NULL,
									mysymbol varchar NOT NULL,
									strikestep float4 NULL DEFAULT 0,
									exchange varchar NULL
								);`
var DB_CREATE_TABLE_USER_SETTING = `CREATE TABLE $1 (
									name varchar NOT NULL,
									controls JSON NOT NULL
								);`

var DB_CREATE_TABLE_USER_STRATEGIES = `CREATE TABLE $1 (
										strategy VARCHAR(100) UNIQUE NOT NULL,
										enabled BOOLEAN NOT NULL DEFAULT 'false',
										engine  VARCHAR(50) NOT NULL,
										trigger_time TIME NOT NULL,
										trigger_days VARCHAR(100) NOT NULL,
										cdl_size SMALLINT NOT NULL,
										instruments TEXT,
										controls JSON
									);`

var DB_CREATE_TABLE_ORDER_BOOK = `CREATE TABLE $1 (
									id SERIAL PRIMARY KEY NOT NULL,
									date DATE NOT NULL,
									instr TEXT NOT NULL,
									strategy  VARCHAR(100) NOT NULL,
									status TEXT,
									instr_id INTEGER,
									dir VARCHAR(50),
									entry DOUBLE PRECISION,
									target DOUBLE PRECISION,
									stoploss DOUBLE PRECISION,
									order_id INTEGER,
									order_trades_entry string,
									order_trade_exit string,
									order_info JSON,
									exit_reason TEXT  DEFAULT 'NA',
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

var DB_VIEW_CREATE = `
					CREATE MATERIALIZED VIEW candles_$1min
					WITH (timescaledb.continuous) AS
					SELECT time_bucket('$1 minutes', time) AS candle, 
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
					`

//  ---------------------------------- QUERIES ----------------------------------
