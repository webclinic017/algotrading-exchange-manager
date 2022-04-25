package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadUserStrategiesFromDb() []appdata.UserStrategies_S {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []appdata.UserStrategies_S

	sqlquery := "SELECT * FROM " + appdata.Env["DB_TBL_USER_STRATEGIES"] + " WHERE enabled = 'true'"

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("user_strategies read error %v\n", err)
		return nil
	}

	return ts
}
