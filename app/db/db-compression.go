package db

import (
	"algo-ex-mgr/app/srv"
	"context"
	"strings"
)

func setupDbCompression(tblName string) {

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	query := strings.ReplaceAll(DB_NSEFUT_COMPRESSION_QUERY, "$1", tblName)

	_, err := myCon.Exec(ctx, query)
	if err != nil {
		srv.WarningLogger.Printf("Error setting up DB Compression: %v\n", err)
	}
}
