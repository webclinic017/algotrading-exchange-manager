package db

import (
	"context"
	"goTicker/app/srv"
)

/*
SELECT tradingsymbol, lot_size
    FROM tracking_symbols ts, instruments i
    WHERE
    		ts.symbol = i.name
    	and
    		mysymbol= 'BANKNIFTY-FUT'
    	and
    		strike >= (0  + 0*ts.strikestep)
		and
    		strike < (0 + ts.strikestep + 0*ts.strikestep)
    	and
    		instrument_type = 'FUT'
    	and
    		expiry > '2022-03-23'
    	and
    		expiry < '2022-04-10'
	ORDER BY
		expiry asc
	LIMIT 10;
*/

func FetchInstrData(instrument string, strikelevel uint64, opdepth int, instrtype string, startdate string, enddate string) (instrname string, lotsize uint16) {

	lock.Lock()
	defer lock.Unlock()
	var size uint16
	var name string

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	sqlQuery := `
	SELECT tradingsymbol, lot_size
		FROM tracking_symbols ts, instruments i
		WHERE 
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

	err := myCon.QueryRow(ctx, sqlQuery,
		instrument, strikelevel, opdepth, instrtype, startdate, enddate).Scan(&name, &size)

	if err != nil {
		srv.ErrorLogger.Printf("FetchOrderData error %v\n", err)
		return "", 0
	}

	return name, size
}
