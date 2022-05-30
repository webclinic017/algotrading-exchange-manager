package srv

import (
	"algo-ex-mgr/app/appdata"
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
	InitTradeLogger()

	InfoLogger.Printf("\n\n#######################################################################################\n|\tNew Instance of algotrading-exchange-manager v0.5.0 (Rel. Date: 30-May-2022) [02]   \n#######################################################################################\n")
}

func InitTradeLogger() {
	logFile, err := os.OpenFile("app/zfiles/log/Trades "+time.Now().Format(time.RFC822)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	mwtl := io.MultiWriter(os.Stdout)
	if err != nil {
		log.Println(err)
	} else {
		mwtl = io.MultiWriter(logFile, os.Stdout)
	}
	TradesLogger = log.New(mwtl, string(appdata.ColorInfo), log.Ltime)

}

func InitLogger() {
	mw := io.MultiWriter(os.Stdout)
	logFile, err := os.OpenFile("app/zfiles/log/algoExchMgr "+time.Now().Format(time.RFC822)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	} else {
		mw = io.MultiWriter(logFile, os.Stdout)
	}

	InfoLogger = log.New(mw, string(appdata.ColorInfo)+" INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, string(appdata.ColorWarning)+"WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, string(appdata.ColorError)+"ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
