package db

import (
	"context"
	"encoding/json"
	"goTicker/app/data"
	"goTicker/app/srv"
)

func StoreTradeSignalInDb(sigData string) uint16 {
	lock.Lock()
	defer lock.Unlock()

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	sqlTradeSig := `INSERT INTO signals_trading (
		strategy_id,
		s_date,
		s_direction,
		s_target,
		s_stoploss,
		s_instr_token)
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
		tradeSignal[0].Strategy_id,
		tradeSignal[0].S_date,
		tradeSignal[0].S_direction,
		tradeSignal[0].S_target,
		tradeSignal[0].S_stoploss,
		tradeSignal[0].S_instr_token)

	if err != nil {
		srv.ErrorLogger.Printf("Unable to insert data into 'symbol ID' database: %v\n", err)
	}

	rows, err := myCon.Query(ctx, `
		SELECT s_order_id 
		FROM signals_trading 
		WHERE  (
				s_instr_token = $1 
			AND 
				s_date = $2
			AND 
				strategy_id = $3)`,
		tradeSignal[0].S_instr_token,
		tradeSignal[0].S_date,
		tradeSignal[0].Strategy_id)

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
			srv.ErrorLogger.Printf("Error: ", rows.Err())
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
