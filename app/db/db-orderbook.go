package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadTradeExitsFromDb() string {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, err := dbPool.Acquire(ctx)
	if err != nil { // DB connection error
		return ""
	}

	defer myCon.Release()

	var e string
	err = myCon.QueryRow(ctx, dbSqlQuery(DB_TRADEMGR_EXISTS_QUERY)).Scan(&e)

	if err != nil {
		srv.TradesLogger.Printf("trademgr - exit conditions read error --> %v\n", err.Error())
		return ""
	}
	return e

}

func ReadOrderBookFromDb(orderBookId uint16) (status bool, tr *appdata.OrderBook_S) {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, err := dbPool.Acquire(ctx)

	if err != nil { // DB connection error
		return false, nil
	}
	defer myCon.Release()

	var or []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlqueryOrderBookId), orderBookId)

	err = pgxscan.Select(ctx, dbPool, &or, sqlquery)

	if err != nil {
		srv.TradesLogger.Printf("order_trades read error --> %v\n", err.Error())
		return false, nil
	}

	if len(or) == 0 {
		srv.TradesLogger.Printf("order_trades - no orders present in db")
		return false, nil
	}

	return true, or[0]

}

func ReadAllOrderBookFromDb(condition string, status string) []*appdata.OrderBook_S {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, err := dbPool.Acquire(ctx)
	if err != nil { // DB connection error
		return nil
	}
	defer myCon.Release()

	var order []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlqueryAllOrderBookCondition), condition, status)

	err = pgxscan.Select(ctx, dbPool, &order, sqlquery)

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
	myCon, err := dbPool.Acquire(ctx)

	if err != nil {
		return tr.Id
	}
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
			tr.ApiSignalEntr,
			tr.ApiSignalExit,
			tr.Orders_entr,
			tr.Orders_exit,
			tr.Post_analysis,
		)
		if err != nil {
			srv.ErrorLogger.Printf("Order Entry : Unable to create new order for strategy-symbol in DB: %v\n", err)
		}

		var c uint16
		err = myCon.QueryRow(ctx, dbSqlQuery(sqlOrderCount),
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
			tr.ApiSignalEntr,
			tr.ApiSignalExit,
			tr.Orders_entr,
			tr.Orders_exit,
			tr.Post_analysis,
			tr.Id,
		)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to update Order for strategy-symbol in DB: %v\n", err)
		}
		return tr.Id
	}
	return tr.Id
}

func FetchOrderData(orderBookId uint16) []*appdata.OrderBook_S {

	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, err := dbPool.Acquire(ctx)

	if err != nil { // DB connection error
		return nil
	}
	defer myCon.Release()

	var ts []*appdata.OrderBook_S

	sqlquery := fmt.Sprintf(dbSqlQuery(sqlqueryOrderBookId), orderBookId)

	err = pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("FetchOrderData error %v\n", err)
		return nil
	}
	return ts
}
