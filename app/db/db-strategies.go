package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"encoding/json"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadStrategiesFromDb() []*appdata.Strategies {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*appdata.Strategies

	err := pgxscan.Select(ctx, dbPool, &ts, `SELECT * FROM strategies where enabled = 'true'`)

	if err != nil {
		srv.ErrorLogger.Printf("Strategies read error %v\n", err)
		return nil
	}

	for each := range ts {
		err = json.Unmarshal([]byte(ts[each].Controls), &ts[each].CtrlParam)
		if err != nil {
			srv.ErrorLogger.Printf("Strategies read error %v\n", err)
		}
	}

	return ts
}
