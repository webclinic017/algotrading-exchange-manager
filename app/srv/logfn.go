package srv

import (
	"io"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	logFile, err := os.OpenFile("app/log/goTickerLogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(logFile, os.Stdout)

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Println("\n\nNew Instance of algotrading-exchange-manager v0.2.0\n\n")
}
