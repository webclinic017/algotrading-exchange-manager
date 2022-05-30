package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func InitTickStorage() {
	go StoreNseIdxFutsInDb()
	go StoreTicksInDb()
}

func StoreNseIdxFutsInDb() {
	var dbTick []appdata.TickData

	for v := range appdata.ChNseTicks { // read from tick channel
		// fmt.Println("Tick: ", appdata.ChNseTicks)
		dbTick = append(dbTick, v)
		if len(dbTick) > 100 {
			dBwg.Add(1)
			go executeBatch(dbTick, appdata.Env["DB_TBL_TICK_NSEFUT"])
			dbTick = nil
		}
	}
	dBwg.Add(1)
	go executeBatch(dbTick, appdata.Env["DB_TBL_TICK_NSEFUT"])
	dbTick = nil

	// wait for all executeBatch() to finish
	dBwg.Wait()
	dbPool.Close()
}

func StoreTicksInDb() {
	var dbTick []appdata.TickData

	for v := range appdata.ChStkTick { // read from tick channel
		// fmt.Println("Tick: ", appdata.ChStkTick)
		dbTick = append(dbTick, v)
		if len(dbTick) > 100 {
			dBwg.Add(1)
			go executeBatch(dbTick, appdata.Env["DB_TBL_TICK_NSESTK"])
			dbTick = nil
		}
	}
	dBwg.Add(1)
	go executeBatch(dbTick, appdata.Env["DB_TBL_TICK_NSESTK"])
	dbTick = nil

	// wait for all executeBatch() to finish
	dBwg.Wait()
	dbPool.Close()
}

func executeBatch(dataTick []appdata.TickData, tableName string) {
	defer dBwg.Done()
	defer func() {
		if err := recover(); err != nil {
			srv.WarningLogger.Print("DB Not intialised: ", err)
		}
	}()

	batch := &pgx.Batch{}

	queryInsertTimeseriesData := `INSERT INTO %v
	(
		time,
		symbol,
		last_traded_price,
		buy_demand, sell_demand,
		trades_till_now,
		open_interest)
		VALUES
		($1, $2, $3, $4, $5, $6, $7);`

	sqlquery := fmt.Sprintf(queryInsertTimeseriesData, tableName)

	for i := range dataTick {
		var ct appdata.TickData = dataTick[i]

		batch.Queue(sqlquery,
			ct.Timestamp,
			ct.Symbol,
			ct.LastTradedPrice,
			ct.Buy_Demand,
			ct.Sell_Demand,
			ct.TradesTillNow,
			ct.OpenInterest)
	}

	ctx := context.Background()

	myCon, err := dbPool.Acquire(ctx)
	if err != nil {
		srv.ErrorLogger.Printf("DB-Ticks : MAJOR ISSUE : Error acquiring connection (this chunk is skipped): %v", err)
	}
	defer myCon.Release()

	// stat := dbpool.Stat()

	br := myCon.SendBatch(ctx, batch)
	_, err = br.Exec()

	if err != nil {
		ErrCnt++
		srv.WarningLogger.Printf("Unable to execute statement in batch queue %v\n", err)
	}
}

// TODO: Major Issue: if ctx fails, write fails to DB - What happens to data? How to recover?
