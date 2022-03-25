package db

import (
	"context"
	"goTicker/app/srv"
)

func FetchInstrData(instrument string, strikelevel uint64, opdepth int, instrtype string, startdate string, enddate string) (instrname string, lotsize float64) {

	lock.Lock()
	defer lock.Unlock()
	var size int64
	var name string

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	sqlQueryOptn := `SELECT tradingsymbol, lot_size
					FROM tracking_symbols ts, instruments i
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
					FROM tracking_symbols ts, instruments i
					WHERE 
						ts.symbol = i.name 
					and 
						ts.mysymbol = $1 
					and
						i.segment = 'NSE'
					and 
						instrument_type = 'EQ'  
					LIMIT 10;`

	sqlQueryFUT := `SELECT tradingsymbol, lot_size
					FROM tracking_symbols ts, instruments i
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
