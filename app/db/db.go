package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dBwg sync.WaitGroup
var lock sync.Mutex
var ErrCnt int = 0
var dbPool *pgxpool.Pool

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
		myCon2.Exec(ctx, DB_CREATE_QUERY)
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
		if createTable(appdata.Env["DB_TBL_TICK_NSEFUT"]+appdata.Env["DB_TEST_PREFIX"], DB_CREATE_TABLE_TICKER) {
			if createTable(appdata.Env["DB_TBL_TICK_NSESTK"]+appdata.Env["DB_TEST_PREFIX"], DB_CREATE_TABLE_TICKER) {
				if createTable(appdata.Env["DB_TBL_PREFIX_USER_ID"]+appdata.Env["DB_TBL_USER_SYMBOLS"]+appdata.Env["DB_TEST_PREFIX"], DB_CREATE_TABLE_USER_SYMBOLS) {
					if createTable(appdata.Env["DB_TBL_PREFIX_USER_ID"]+appdata.Env["DB_TBL_USER_SETTING"]+appdata.Env["DB_TEST_PREFIX"], DB_CREATE_TABLE_USER_SETTING) {
						if createTable(appdata.Env["DB_TBL_PREFIX_USER_ID"]+appdata.Env["DB_TBL_USER_STRATEGIES"]+appdata.Env["DB_TEST_PREFIX"], DB_CREATE_TABLE_USER_STRATEGIES) {
							if createTable(appdata.Env["DB_TBL_PREFIX_USER_ID"]+appdata.Env["DB_TBL_ORDER_BOOK"]+appdata.Env["DB_TEST_PREFIX"], DB_CREATE_TABLE_ORDER_BOOK) {
								// createViews()
								setupDbCompression(appdata.Env["DB_TICK_TABLE_NSEFUT"] + appdata.Env["DB_TEST_PREFIX"])
								srv.InfoLogger.Printf("DB checks completed\n")
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func CloseDb() {
	dbPool.Close()
}

func DbRawExec(query string) {
	ctx := context.Background()
	myCon, err := dbPool.Acquire(ctx)
	defer myCon.Release()
	if err != nil {
		srv.ErrorLogger.Printf("Could not acquire Context, too many operations?: %v\n", err)
		return
	}

	myCon.Exec(ctx, query)
}
