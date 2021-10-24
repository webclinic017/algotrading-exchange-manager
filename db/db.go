package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

func DbInit() bool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	ctx := context.Background()
	var err error

	dbPool, err = pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return false
	}
	defer dbPool.Close()

	var greeting string
	err = dbPool.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return false
	}
	fmt.Println("connected to DB : " + greeting)

	// check if table exist, else create it
	queryCreateTicksTable := `CREATE TABLE ticks (
                                                time TIMESTAMPTZ NOT NULL,
                                                symbol integer NULL,
                                                last_price double precision NULL,
                                                open double precision NULL,
                                                close double precision NULL,
                                                low double precision NULL,
                                                high double precision NULL,
                                                volume int NULL
                                            );
						SELECT create_hypertable('ticks', 'time');
						`

	//execute statement, fails if table already exists
	_, _ = dbPool.Exec(ctx, queryCreateTicksTable)

	return true

}

func StoreTickInDb() bool {

	return true
}
