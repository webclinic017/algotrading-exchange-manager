package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadOrderBookFromDb(orderBookId uint16) (status bool, tr *appdata.OrderBook_S) {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlqueryOrderBookId), orderBookId)

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err.Error())
		return false, nil

	}

	if len(ts) == 0 {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err)
		return false, nil
	}

	return true, ts[0]

}

func ReadAllActiveOrderBookFromDb() []*appdata.OrderBook_S {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlQueryAllActiveOrderBook), "TradeCompleted")

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

func ReadAllOrderBookFromDb(condition string, status string) []*appdata.OrderBook_S {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var order []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlqueryAllOrderBookCondition), condition, status)

	err := pgxscan.Select(ctx, dbPool, &order, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("order_trades read error %v\n", err.Error())
		return nil

	}

	if len(order) == 0 {
		srv.InfoLogger.Printf("order_trades 0 %v\n", err)
		return nil
	}

	return order

}

func StoreOrderBookInDb(tr appdata.OrderBook_S) uint16 {
	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	if tr.Id == 0 {

		_, err := myCon.Exec(ctx, dbSqlQuery(sqlCreateOrder),
			tr.Date,
			tr.Instr,
			tr.Strategy,
			tr.Status,
			tr.Dir,
			tr.Exit_reason,
			tr.Info,
			tr.Targets,
			tr.Orders_entr,
			tr.Orders_exit,
			tr.Post_analysis,
		)
		if err != nil {
			srv.ErrorLogger.Printf("Order Entry : Unable to create new order for strategy-symbol in DB: %v\n", err)
		}
	} else {

		_, err := myCon.Exec(ctx, dbSqlQuery(sqlUpdateOrder),
			tr.Date,
			tr.Instr,
			tr.Strategy,
			tr.Status,
			tr.Dir,
			tr.Exit_reason,
			tr.Info,
			tr.Targets,
			tr.Orders_entr,
			tr.Orders_exit,
			tr.Post_analysis,
			tr.Id,
		)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to update Order for strategy-symbol in DB: %v\n", err)
		}
	}

	var c uint16
	err := myCon.QueryRow(ctx, dbSqlQuery(sqlOrderCount),
		tr.Instr,
		tr.Date,
		tr.Strategy).Scan(&c)
	// RULE: Instrument, Date, Strategy (combined) must be unique

	if err != nil {
		srv.ErrorLogger.Printf("OrderBook DB store error %v\n", err)
		return 0
	}

	if c == 1 {
		err = myCon.QueryRow(ctx, dbSqlQuery(sqlOrderId),
			tr.Instr,
			tr.Date,
			tr.Strategy).Scan(&c)

		if err != nil {
			srv.ErrorLogger.Printf("OrderBook DB store error %v\n", err)
			return 0
		}
		return c
	}
	return 0
}

func FetchOrderData(orderBookId uint16) []*appdata.OrderBook_S {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlqueryOrderBookId), orderBookId)

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("FetchOrderData error %v\n", err)
		return nil
	}

	return ts

}

// t_entry = 0
