package db

import (
	"context"
	"goTicker/app/data"

	"github.com/georgysavva/scany/pgxscan"
)

func ReadStrategiesFromDb() []*data.Strategies {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	var ts []*data.Strategies

	pgxscan.Select(ctx, dbPool, &ts, `SELECT * FROM strategies where strategy_en = 'true'`)

	return ts
}
