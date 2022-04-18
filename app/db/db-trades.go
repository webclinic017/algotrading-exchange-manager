package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadTradeSignalFromDb(orderBookId uint16) (status bool, tr *appdata.TradeSignal) {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.TradeSignal

	sqlquery := fmt.Sprintf("SELECT * FROM order_trades WHERE id = %d", orderBookId)

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err)
		return false, nil

	}

	if len(ts) == 0 {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err)
		return false, nil
	}

	return true, ts[0]

}

func ReadAllTradeSignalFromDb() []*appdata.TradeSignal {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.TradeSignal

	sqlquery := fmt.Sprintf("SELECT * FROM order_trades WHERE status != '%s'", "TradeCompleted")

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err)
		return nil

	}

	if len(ts) == 0 {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err)
		return nil
	}

	return ts

}

func StoreTradeSignalInDb(tr appdata.TradeSignal) uint16 {
	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	// if sigData != "" { // signal found, parse json
	// 	var apiSignal []*appdata.ApiSignal
	// 	err := json.Unmarshal([]byte(sigData), &apiSignal)
	// 	if err != nil {
	// 		srv.TradesLogger.Printf("apiSignal - API JSON data parse error: %v\n", err)
	// 	}
	// 	tr.Dir = apiSignal[0].Dir
	// 	tr.Entry = apiSignal[0].Entry
	// 	tr.Target = apiSignal[0].Target
	// 	tr.Stoploss = apiSignal[0].Stoploss
	// }

	var sqlquery string
	sqlCreateTradeSig := `INSERT INTO order_trades (
		date,
		instr,
		strategy,
		status,
		instr_id,
		dir,
		entry,
		target,
		stoploss,
		order_id,
		order_trades_entry,
		order_trade_exit,
		order_simulation,
		exit_reason,
		post_analysis)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`

	sqlUpdateTradeSig := `INSERT INTO order_trades (
			date,
			instr,
			strategy,
			status,
			instr_id,
			dir,
			entry,
			target,
			stoploss,
			order_id,
			order_trades_entry,
			order_trade_exit,
			order_simulation,
			exit_reason,
			post_analysis)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`

	if tr.Id == 0 {
		sqlquery = sqlCreateTradeSig
	} else {
		sqlquery = sqlUpdateTradeSig
	}

	_, err := myCon.Exec(ctx, sqlquery,
		tr.Date,
		tr.Instr,
		tr.Strategy,
		tr.Status,
		tr.Instr_id,
		tr.Dir,
		tr.Entry,
		tr.Target,
		tr.Stoploss,
		tr.Order_id,
		tr.Order_trades_entry,
		tr.Order_trades_exit,
		tr.Order_simulation,
		tr.Exit_reason,
		tr.Post_analysis,
	)

	if err != nil {
		srv.ErrorLogger.Printf("Unable to insert strategy-symbol in DB: %v\n", err)
	}

	rows, err := myCon.Query(ctx, `
		SELECT id 
		FROM order_trades 
		WHERE  (
				instr = $1 
			AND 
				date = $2
			AND 
				strategy = $3)`,
		tr.Instr,
		tr.Date,
		tr.Strategy)
	// RULE: Instrument, Date, Strategy (combined) must be unique

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
			srv.ErrorLogger.Printf("TradeSignal DB row-scan error %v\n", err)
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
	} else if (len(orderId)) > 1 {
		srv.ErrorLogger.Printf("TradeSignal - Multiple entries in DB - Skipping trades for %v %v\n", tr.Strategy, err)
	} else {
		srv.ErrorLogger.Printf("TradeSignal DB unkown error %v\n", err)

	}
	return 0
}

func FetchOrderData(orderBookId uint16) []*appdata.TradeSignal {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.TradeSignal

	sqlquery := fmt.Sprintf("SELECT * FROM signals_trading WHERE id = %d", orderBookId)

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("FetchOrderData error %v\n", err)
		return nil
	}

	return ts

}

// t_entry = 0
