package db

import (
	"context"
	"encoding/json"
	"goTicker/app/data"
	"goTicker/app/srv"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadStrategiesFromDb() []*data.Strategies {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*data.Strategies

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
