package db

var DB_EXISTS_QUERY = "SELECT datname FROM pg_catalog.pg_database  WHERE lower(datname) = lower('algotrading');"
var DB_CREATE_QUERY = "CREATE DATABASE algotrading;"
var DB_CREATE_TABLE_ID_DECODED = `CREATE TABLE token_id_decoded
								(
									time TIMESTAMP NOT NULL,
									nse_symbol VARCHAR(30),
									mcx_symbol VARCHAR(30)
								);
						`
var DB_TABLE_TICKER_NAME_NSE_IDX = `ticks_nsefut`
var DB_TABLE_TICKER_NAME_STK = `ticks_stk`
var DB_CREATE_TABLE_NSE_IDX_TICKER = `CREATE TABLE ticks_nsefut
							 		(
										time TIMESTAMP NOT NULL,
										symbol VARCHAR(30) NOT NULL,
										last_traded_price double precision NOT NULL DEFAULT 0,
										buy_demand bigint NOT NULL DEFAULT 0,
										sell_demand bigint NOT NULL DEFAULT 0,
										trades_till_now bigint NOT NULL DEFAULT 0,
										open_interest bigint NOT NULL DEFAULT 0
									);
								SELECT create_hypertable('ticks_nsefut', 'time');
								SELECT set_chunk_time_interval('ticks_nsefut', INTERVAL '7 days');
						`
var DB_CREATE_TABLE_STK_TICKER = `CREATE TABLE ticks_stk
									(
										time TIMESTAMP NOT NULL,
										symbol VARCHAR(30) NOT NULL,
										last_traded_price double precision NOT NULL DEFAULT 0,
										buy_demand bigint NOT NULL DEFAULT 0,
										sell_demand bigint NOT NULL DEFAULT 0,
										trades_till_now bigint NOT NULL DEFAULT 0,
										open_interest bigint NOT NULL DEFAULT 0
									);
									SELECT create_hypertable('ticks_stk', 'time');
									SELECT set_chunk_time_interval('ticks_stk', INTERVAL '7 days');`

var DB_NSEFUT_COMPRESSION_QUERY = `ALTER TABLE ticks_nsefut SET 
							(
								timescaledb.compress,
								timescaledb.compress_segmentby = 'symbol'
							); 
							SELECT add_compression_policy('ticks_nsefut ', INTERVAL '30 days');
						`

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
