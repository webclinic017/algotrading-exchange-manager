package db

import (
	"context"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"os"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dBwg sync.WaitGroup
var lock sync.Mutex
var ErrCnt int = 0
var dbPool *pgxpool.Pool
var dbTick []kite.TickData

func connectDB() bool {
	ctx := context.Background()
	dbUrl := "postgres://" + os.Getenv("TIMESCALEDB_USERNAME") + ":" + os.Getenv("TIMESCALEDB_PASSWORD") + "@" + os.Getenv("TIMESCALEDB_ADDRESS") + ":" + os.Getenv("TIMESCALEDB_PORT") + "/postgres"

	// Check if you can connect to DB server (accessing 'postgres' defualt DB)
	dbPoolDefault, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		srv.ErrorLogger.Println("Could not connect with 'postgres' DB\n", err)
		return false
	}
	defer dbPoolDefault.Close()

	myCon, err := dbPoolDefault.Acquire(ctx)
	defer myCon.Release()
	if err != nil {
		srv.ErrorLogger.Printf("Could not acquire Context, too many operations?: %v\n", err)
		return false
	}

	// check if 'algotrading' database exists, if not CREATE it
	var retVal string
	myCon.QueryRow(ctx, DB_EXISTS_QUERY).Scan(&retVal)

	if len(retVal) == 0 {
		srv.InfoLogger.Printf("algotrading DB Does not exist, creating now!: %v\n", err)

		//execute statement, fails if table already exists
		myCon2, _ := dbPoolDefault.Acquire(ctx)
		defer myCon.Release()
		_, err = myCon2.Exec(ctx, "DB_CREApackage dbd to CREATE algotrading DB: %v\n", err)
		return false
	}
	return true
}

func DbInit() bool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	srv.InfoLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Db Checks",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	ctx := context.Background()
	dbUrl := "postgres://" + os.Getenv("TIMESCALEDB_USERNAME") + ":" + os.Getenv("TIMESCALEDB_PASSWORD") + "@" + os.Getenv("TIMESCALEDB_ADDRESS") + ":" + os.Getenv("TIMESCALEDB_PORT") + "/algotrading"

	if connectDB() {
		// 1. Connect with 'algotrading' DB
		var err error
		dbPool, err = pgxpool.Connect(ctx, dbUrl)
		if err != nil {
			srv.ErrorLogger.Printf("Unable to connect with 'algotrading db' %v\n", err)
			return false
		}
		// 2. Aquire context
		myCon, err := dbPool.Acquire(ctx)
		defer myCon.Release()
		if err != nil {
			srv.ErrorLogger.Printf("Could not acquire Context, too many operations?: %v\n", err)
			return false
		}

		// 3. Check if 'ticker' table exists, if not CREATE it
		if createTable(DB_TABLE_TICKER_NAME, DB_CREATE_TABLE_TICKER) {
			createViews()
			setupDbCompression()
			createTable(DB_TABLE_ID_DECODED_NAME, DB_CREATE_TABLE_ID_DECODED)
			srv.InfoLogger.Printf("DB checks completed\n")
			return true
		} else {
			return false
		}

	} else {
		return false
	}

}

func CloseDb() {
	dbPool.Close()
}
