package db

import (
	"context"
	"goTicker/app/srv"
	"strings"

	"github.com/jackc/pgconn"
)

func createViews() {

	createViewInMinutes("1")
	createViewInMinutes("3")
	createViewInMinutes("5")
	createViewInMinutes("10")
	createViewInMinutes("15")

}

func createViewInMinutes(viewMin string) {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	if !viewExists("candles_" + viewMin + "min") {

		sqlquery := strings.Replace(DB_CREATE_VIEW, "$1", viewMin, -1)

		// $1 candles_5min // $2 5
		_, err := myCon.Exec(ctx, sqlquery)
		if err != nil {
			pgerr, _ := err.(*pgconn.PgError)
			if pgerr.Code != "42P07" {
				srv.WarningLogger.Printf("Error creating candles_"+viewMin+"min: %v\n", err)
			}
		}
	}
}

func viewExists(viewName string) bool {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var retVal string

	// query := "SELECT view_name FROM timescaledb_information.continuous_aggregates  WHERE view_name = '" + viewName + "';"
	err := myCon.QueryRow(ctx, DB_VIEW_EXISTS_QUERY, viewName).Scan(&retVal)
	if err != nil {
		println(err.Error())
	}

	if len(retVal) == 0 {
		return false
	}

	return true
}
