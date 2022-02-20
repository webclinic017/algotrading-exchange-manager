package srv

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	TradesLogger  *log.Logger
)

// TODO: update logger module for date wise operations

func Init() {
	InitLogger()
	InfoLogger.Printf("\n\n###########################################################################\n\tNew Instance of algotrading-exchange-manager v0.4.0\n###########################################################################\n")
}

func InitTradeLogger() {
	logFile, err := os.OpenFile("app/log/Trades "+time.Now().Format(time.RFC822)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	TradesLogger = log.New(logFile, "", log.Ltime)
}

func InitLogger() {
	logFile, err := os.OpenFile("app/log/algoExchMgr "+time.Now().Format(time.RFC822)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(logFile, os.Stdout)

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
