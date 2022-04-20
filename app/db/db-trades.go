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

	tblName := appdata.Env["DB_TBL_PREFIX_USER_ID"] +
		appdata.Env["DB_TBL_ORDER_BOOK"] +
		appdata.Env["DB_TEST_PREFIX"]

	sqlquery := fmt.Sprintf("SELECT * FROM "+tblName+" WHERE id = %d", orderBookId)

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

func ReadAllActiveTradeSignalFromDb() []*appdata.TradeSignal {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.TradeSignal

	tblName := appdata.Env["DB_TBL_PREFIX_USER_ID"] +
		appdata.Env["DB_TBL_ORDER_BOOK"] +
		appdata.Env["DB_TEST_PREFIX"]
	sqlquery := fmt.Sprintf("SELECT * FROM "+tblName+" WHERE status ! = %d", "TradeCompleted")

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

func ReadAllTradeSignalFromDb(condition string, status string) []*appdata.TradeSignal {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.TradeSignal

	tblName := appdata.Env["DB_TBL_PREFIX_USER_ID"] +
		appdata.Env["DB_TBL_ORDER_BOOK"] +
		appdata.Env["DB_TEST_PREFIX"]

	sqlquery := fmt.Sprintf("SELECT * FROM " + tblName + " WHERE status " + condition + " '" + status + "'")

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

	tblName := appdata.Env["DB_TBL_PREFIX_USER_ID"] +
		appdata.Env["DB_TBL_ORDER_BOOK"] +
		appdata.Env["DB_TEST_PREFIX"]

	if tr.Id == 0 {
		sqlCreateTradeSig := `INSERT INTO ` + tblName + `(
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
			order_trade_entry,
			order_trade_exit,
			order_simulation,
			exit_reason,
			post_analysis)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`

		_, err := myCon.Exec(ctx, sqlCreateTradeSig,
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
			tr.Order_trade_entry,
			tr.Order_trade_exit,
			tr.Order_simulation,
			tr.Exit_reason,
			tr.Post_analysis,
		)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to insert strategy-symbol in DB: %v\n", err)
		}
	} else {

		sqlUpdateTradeSig := ` UPDATE ` + tblName + ` SET 
			date = $1,
			instr = $2,
			strategy = $3,
			status = $4,
			instr_id = $5,
			dir = $6,
			entry = $7,
			target = $8,
			stoploss = $9,
			order_id = $10,
			order_trade_entry = $11,
			order_trade_exit = $12,
			order_simulation = $13,
			exit_reason = $14,
			post_analysis = $15
			WHERE id = $16;`

		_, err := myCon.Exec(ctx, sqlUpdateTradeSig,
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
			tr.Order_trade_entry,
			tr.Order_trade_exit,
			tr.Order_simulation,
			tr.Exit_reason,
			tr.Post_analysis,
			tr.Id,
		)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to insert strategy-symbol in DB: %v\n", err)
		}
	}

	sqquery := `
SELECT id 
FROM ` + tblName + ` 
WHERE  (
		instr = $1 
	AND 
		date = $2
	AND 
		strategy = $3)`

	rows, err := myCon.Query(ctx, sqquery,
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
