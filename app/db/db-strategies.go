package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"encoding/json"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadStrategiesFromDb() []appdata.Strategies {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []appdata.Strategies

	tblName := appdata.Env["DB_TBL_PREFIX_USER_ID"] + appdata.Env["DB_TBL_USER_STRATEGIES"] + appdata.Env["DB_TEST_PREFIX"]
	sqlquery := "SELECT * FROM " + tblName + " WHERE enabled = 'true'"

	err := pgxscan.Select(ctx, dbPool, &ts, sqlquery)

	if err != nil {
		srv.ErrorLogger.Printf("user_strategies read error %v\n", err)
		return nil
	}

	for each := range ts {
		err = json.Unmarshal([]byte(ts[each].Controls), &ts[each].CtrlParam)
		if err != nil {
			srv.ErrorLogger.Printf("user_strategies read error %v\n", err)
		}
	}

	return ts
}
