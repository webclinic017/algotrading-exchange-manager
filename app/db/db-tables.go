package db

import (
	"context"
	"goTicker/app/srv"
)

func createTable(tblName string, sqlquery string) bool {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)

	var retVal string

	query := "select table_name from information_schema.tables WHERE table_name = '" + tblName + "';"
	myCon.QueryRow(ctx, query).Scan(&retVal)
	// if err != nil {
	// 	srv.WarningLogger.Printf("Failed to CREATE %s table : %v\n", tblName, err)
	// }

	if len(retVal) == 0 {
		srv.InfoLogger.Printf("%s Does not exist, creating now!\n", tblName)
		_, err := myCon.Exec(ctx, sqlquery)
		if err != nil {
			srv.WarningLogger.Printf("Failed to CREATE %s table : %v\n", tblName, err)
			myCon.Release()
			return false
		}
	}
	myCon.Release()
	return true
}
