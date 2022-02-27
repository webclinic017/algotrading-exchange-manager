package db

import (
	"context"
	"goTicker/app/srv"
	"time"
)

func StoreSymbolsInDb(nse_symbol string, mcx_symbol string) {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	timestamp := time.Now()
	queryInsertMetadata := `INSERT INTO token_id_decoded (
		time,
		nse_symbol,
		mcx_symbol)
		VALUES
		($1, $2, $3);`

	_, err := myCon.Exec(ctx, queryInsertMetadata,
		timestamp,
		nse_symbol,
		mcx_symbol)
	if err != nil {
		srv.ErrorLogger.Printf("Unable to insert data into 'symbol ID' database: %v\n", err)
	}
}
