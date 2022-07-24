package db

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dBwg sync.WaitGroup
var lock sync.Mutex
var ErrCnt int = 0
var dbPool *pgxpool.Pool

func connectDB() bool {
	ctx := context.Background()

	// Check if you can connect to DB server (accessing 'postgres' defualt DB)
	dbPoolDefault, err := pgxpool.Connect(context.Background(), appdata.Env["DB_URL"]+"/postgres")
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

	if connectDB() {
		// 1. Connect with 'algotrading' DB
		var err error
		dbPool, err = pgxpool.Connect(ctx, appdata.Env["DB_URL"]+"/algotrading")
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

		var s bool
		s = createTable("DB_TBL_TICK_NSEFUT", DB_CREATE_TABLE_TICKER_NSEFUT)
		s = s && createTable("DB_TBL_TICK_NSESTK", DB_CREATE_TABLE_TICKER_NSESTK)
		s = s && createTable("DB_TBL_USER_SYMBOLS", DB_CREATE_TABLE_USER_SYMBOLS)
		s = s && createTable("DB_TBL_USER_SETTING", DB_CREATE_TABLE_USER_SETTING)
		s = s && createTable("DB_TBL_USER_STRATEGIES", DB_CREATE_TABLE_USER_STRATEGIES)
		s = s && createTable("DB_TBL_ORDER_BOOK", DB_CREATE_TABLE_ORDER_BOOK)
		s = s && createTable("DB_TBL_CDL_VIEW_STK", DB_VIEW_CREATE_STK)
		s = s && createTable("DB_TBL_CDL_VIEW_FUT", DB_VIEW_CREATE_FUT)

		// schedule views to run @ 5pm everyday
		views_reschedule()

		if s {
			// setupDbCompression(appdata.Env["DB_TICK_TABLE_NSEFUT"])
			srv.InfoLogger.Printf("DB checks completed\n")
			return true
		}

	}
	return false
}

func views_reschedule() {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)

	var retVal int

	dt := time.Now().Format("2006-01-02")
	dt = dt + " 17:00:00.000 +0530"

	myCon.QueryRow(ctx, dbSqlQuery(sqlQueryViewStkGetID)).Scan(&retVal)
	query := "SELECT alter_job(" + fmt.Sprintf("%v", retVal) + ", next_start => '" + dt + "');"
	myCon.Exec(ctx, query)

	myCon.QueryRow(ctx, dbSqlQuery(sqlQueryViewFutGetID)).Scan(&retVal)
	query = "SELECT alter_job(" + fmt.Sprintf("%v", retVal) + ", next_start => '" + dt + "');"
	myCon.Exec(ctx, query)

	myCon.Release()

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
	_, err = myCon.Exec(ctx, query)
	if err != nil {
		srv.ErrorLogger.Printf("Error executing query: %v\n", err.Error())
	}

}

func createTable(tblName string, sqlquery string) bool {
	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)

	var retVal string

	query := "select table_name from information_schema.tables WHERE table_name = '" + appdata.Env[tblName] + "';"
	myCon.QueryRow(ctx, query).Scan(&retVal)

	if len(retVal) == 0 {
		srv.InfoLogger.Printf("%s{%s} Does not exist, creating now!\n", tblName, appdata.Env[tblName])

		_, err := myCon.Exec(ctx, dbSqlQuery(sqlquery))
		if err != nil {
			srv.WarningLogger.Printf("Failed to CREATE %s table : %v\n", tblName, err.Error())
			myCon.Release()
			return false
		}
	}
	myCon.Release()
	return true
}

func setupDbCompression(tblName string) {

	ctx := context.Background()
	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	query := strings.ReplaceAll(DB_NSEFUT_COMPRESSION_QUERY, "$1", tblName)

	_, err := myCon.Exec(ctx, query)
	if err != nil {
		srv.WarningLogger.Printf("Error setting up DB Compression: %v\n", err)
	}
}

func DbSaveInstrCsv(table string, filePath string) {

	defer func() {
		if err := recover(); err != nil {
			srv.WarningLogger.Print("DB Not intialised: ", err)
		}
	}()

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	batch := &pgx.Batch{}

	for _, record := range records {

		if table == "instruments" {
			batch.Queue(dbSqlQuery(sqlSaveInstruments), record[0], record[1], record[2], record[3],
				record[4], record[5], record[6], record[7],
				record[8], record[9], record[10], record[11])
		} else if table == "user_symbols" {
			batch.Queue(dbSqlQuery(sqlSaveUserSymbols),
				record[0], record[1], record[2],
				record[3], record[4], record[5])
		}
	}

	ctx := context.Background()

	myCon, _ := dbPool.Acquire(ctx)
	defer myCon.Release()

	if table == "instruments" {
		myCon.Exec(ctx, dbSqlQuery(DB_CREATE_TBL_INSTRUMENTS))
	} else if table == "user_symbols" {
		myCon.Exec(ctx, dbSqlQuery(DB_CREATE_TABLE_USER_SYMBOLSwDel))
	}

	br := myCon.SendBatch(ctx, batch)
	_, err = br.Exec()

	// fmt.Println("Inserted ", ct, " rows")
	if err != nil {
		ErrCnt++
		srv.WarningLogger.Printf("Unable to execute statement in batch queue %v\n", err)
	}

	if table == "instruments" {
		time.Sleep(time.Second * 5)
	}
}
