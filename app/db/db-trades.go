package db

import (
	"context"
	"encoding/json"
	"fmt"
	"goTicker/app/data"
	"goTicker/app/srv"

	"github.com/georgysavva/scany/pgxscan"
)

func StoreTradeSignalInDb(sigData string) uint16 {
	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	sqlTradeSig := `INSERT INTO signals_trading (
		strategy,
		date,
		dir,
		target,
		stoploss,
		instr)
		VALUES
		($1, $2, $3, $4, $5, $6);`

	var tradeSignal []*data.TradeSignal
	err := json.Unmarshal([]byte(sigData), &tradeSignal)
	if err != nil {
		srv.ErrorLogger.Printf("TradeSignal - API JSON data parse error: %v\n", err)
	}

	// // fmt.Printf("%+v\n", tradeSignal[0])
	// // fmt.Println(tradeSignal[0].T_entry)

	_, err = myCon.Exec(ctx, sqlTradeSig,
		tradeSignal[0].Strategy,
		tradeSignal[0].Date,
		tradeSignal[0].Dir,
		tradeSignal[0].Target,
		tradeSignal[0].Stoploss,
		tradeSignal[0].Instr)

	if err != nil {
		srv.ErrorLogger.Printf("Unable to insert data into 'symbol ID' database: %v\n", err)
	}

	rows, err := myCon.Query(ctx, `
		SELECT id 
		FROM signals_trading 
		WHERE  (
				instr = $1 
			AND 
				date = $2
			AND 
				strategy = $3)`,
		tradeSignal[0].Instr,
		tradeSignal[0].Date,
		tradeSignal[0].Strategy)

	if err != nil {
		srv.ErrorLogger.Printf("TradeSignal DB store error %v\n", err)
		return 0
	}

	var orderId []uint16
	// var err1 error

	for rows.Next() {

		var id uint16
		err = rows.Scan(&id)
		if err != nil {
			srv.ErrorLogger.Printf("TradeSignal DB store error %v\n", err)
			rows.Close()
			return 0
		}
		orderId = append(orderId, id)

		if rows.Err() != nil {
			srv.ErrorLogger.Println("Error: ", rows.Err())
			rows.Close()
			return 0
		}
	}
	rows.Close()

	if (len(orderId)) == 1 {
		return orderId[0]
	} else {
		return 0
	}

}

func FetchOrderData(orderBookId uint16) []*data.TradeSignal {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*data.TradeSignal

	sqlquery := fmt.Sprintf("SELECT * FROM signals_trading WHERE id = %d", orderBookId)

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("TradeSignal DB store error %v\n", err)
		return nil
	}

	return ts

}

// t_entry = 0
