package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"strconv"
)

func GetInstrumentsToken() map[string]string {

	var tokensMap = make(map[string]string)

	tknEq := getNseEqTokens()
	tknFut := getFuturesTokens()

	for k, v := range tknEq {
		tokensMap[k] = v
	}
	for k, v := range tknFut {
		tokensMap[k] = v
	}
	return tokensMap
}

func getFuturesTokens() map[string]string {

	var tokensMap = make(map[string]string)

	sqlQueryFutures := `SELECT i.instrument_token, ts.mysymbol
    FROM user_symbols ts, tick_instr i
    WHERE 
    		ts.symbol = i.name
		and 
			ts.segment = i.instrument_type 
		and 
			ts.exchange = i.exchange
    	and 
    		EXTRACT(MONTH FROM TO_DATE(i.expiry,'YYYY-MM-DD')) = EXTRACT(MONTH FROM current_date);`

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	rows, err := myCon.Query(ctx, sqlQueryFutures)

	if err != nil {
		srv.ErrorLogger.Printf("Cannot read list of tokens for ticker %v\n", err)
		return tokensMap
	}

	for rows.Next() {

		var itoken int64
		var symbol string

		err = rows.Scan(&itoken, &symbol)
		if err != nil {
			srv.ErrorLogger.Printf("Cannot parse list of tokens for ticker %v\n", err)
			return tokensMap
		}

		if rows.Err() != nil {
			srv.ErrorLogger.Println("Cannot parse list of tokens for ticker: ", rows.Err())

			return tokensMap
		}
		tokensMap[strconv.FormatInt(itoken, 10)] = symbol
	}
	defer rows.Close()

	return tokensMap
}

// TODO: Check for MCX, BSE and other symbols
func getNseEqTokens() map[string]string {

	var tokensMap = make(map[string]string)

	sqlQueryStocks := `SELECT i.instrument_token, ts.mysymbol
	FROM user_symbols ts, tick_instr i
	WHERE 
			ts.symbol = i.tradingsymbol
		and 
			i.instrument_type = 'EQ'
		and 
			ts.exchange = i.exchange;`

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	rows, err := myCon.Query(ctx, sqlQueryStocks)

	if err != nil {
		srv.ErrorLogger.Printf("Cannot read list of tokens for ticker %v\n", err)
		return tokensMap
	}

	for rows.Next() {

		var itoken int64
		var symbol string

		err = rows.Scan(&itoken, &symbol)
		if err != nil {
			srv.ErrorLogger.Printf("Cannot parse list of tokens for ticker %v\n", err)
			return tokensMap
		}

		if rows.Err() != nil {
			srv.ErrorLogger.Println("Cannot parse list of tokens for ticker: ", rows.Err())

			return tokensMap
		}
		tokensMap[strconv.FormatInt(itoken, 10)] = symbol
	}
	defer rows.Close()

	return tokensMap
}

func FetchInstrData(instrument string, strikelevel uint64, opdepth int, instrtype string, startdate string, enddate string) (instrname string, lotsize float64) {

	lock.Lock()
	defer lock.Unlock()
	var size int64
	var name string

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	tblName := appdata.Env["DB_TBL_PREFIX_USER_ID"] +
		appdata.Env["DB_TBL_USER_SYMBOLS"] +
		appdata.Env["DB_TEST_PREFIX"]

	sqlQueryOptn := `SELECT tradingsymbol, lot_size
					FROM ` + tblName + ` ts, ticks_instr i
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

	sqlQueryEQ := `SELECT tradingsymbol, lot_size
					FROM ` + tblName + ` ts, ticks_instr i
					WHERE 
						ts.symbol = i.tradingsymbol 
					and 
						ts.mysymbol = $1 
					and
						i.segment = 'NSE'
					and 
						instrument_type = 'EQ'  
					LIMIT 10;`

	sqlQueryFUT := `SELECT tradingsymbol, lot_size
					FROM ` + tblName + ` ts, ticks_instr i
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

	var err error
	if instrtype == "EQ" {

		err = myCon.QueryRow(ctx, sqlQueryEQ,
			instrument).Scan(&name, &size)

	} else if instrtype == "FUT" {

		err = myCon.QueryRow(ctx, sqlQueryFUT, instrument,
			startdate, enddate).Scan(&name, &size)

	} else {

		err = myCon.QueryRow(ctx, sqlQueryOptn,
			instrument, strikelevel,
			opdepth, instrtype,
			startdate, enddate).Scan(&name, &size)
	}

	if err != nil {
		srv.ErrorLogger.Printf("FetchOrderData error %v\n", err)
		return "", 0
	}

	return name, float64(size)
}
