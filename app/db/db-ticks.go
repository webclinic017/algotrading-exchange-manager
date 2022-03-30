package db

import (
	"algo-ex-mgr/app/kite"
	"algo-ex-mgr/app/srv"
	"context"

	"github.com/jackc/pgx/v4"
)

func StoreTickInDb() {
	for v := range kite.ChTick { // read from tick channel
		dbTick = append(dbTick, v)
		if len(dbTick) > 250 {
			dBwg.Add(1)
			go executeBatch(dbTick)
			dbTick = nil
		}
	}
	dBwg.Add(1)
	go executeBatch(dbTick)
	dbTick = nil

	// wait for all executeBatch() to finish
	dBwg.Wait()
	dbPool.Close()
}

func executeBatch(dataTick []kite.TickData) {
	defer dBwg.Done()
	defer func() {
		if err := recover(); err != nil {
			srv.WarningLogger.Print("DB Not intialised: ", err)
		}
	}()

	batch := &pgx.Batch{}

	queryInsertTimeseriesData := `INSERT INTO ticks_data
 (
					time,
					symbol,
					last_traded_price,
					buy_demand, sell_demand,
					trades_till_now,
					open_interest)
					VALUES
					($1, $2, $3, $4, $5, $6, $7);`

	for i := range dataTick {
		var ct kite.TickData = dataTick[i]

		batch.Queue(queryInsertTimeseriesData,
			ct.Timestamp,
			ct.Symbol,
			ct.LastTradedPrice,
			ct.Buy_Demand,
			ct.Sell_Demand,
			ct.TradesTillNow,
			ct.OpenInterest)
	}

	ctx := context.Background()

	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	// stat := dbpool.Stat()

	br := myCon.SendBatch(ctx, batch)
	_, err := br.Exec()

	if err != nil {
		ErrCnt++
		srv.WarningLogger.Printf("Unable to execute statement in batch queue %v\n", err)
	}
}
