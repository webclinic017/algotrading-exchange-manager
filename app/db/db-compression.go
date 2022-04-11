package db

import (
	"algo-ex-mgr/app/srv"
	"context"
)

func setupDbCompression() {

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	_, err := myCon.Exec(ctx, DB_NSEFUT_COMPRESSION_QUERY)
	if err != nil {
		srv.WarningLogger.Printf("Error setting up DB Compression: %v\n", err)
	}
}
