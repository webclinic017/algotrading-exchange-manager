package srv

import (
	"fmt"
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

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// TODO: update logger module for date wise operations

func Init() {
	InitLogger()

	fmt.Printf(InfoColor, "Info")
	fmt.Println("")
	fmt.Printf(NoticeColor, "Notice")
	fmt.Println("")
	fmt.Printf(WarningColor, "Warning")
	fmt.Println("")
	fmt.Printf(ErrorColor, "Error")
	fmt.Println("")
	fmt.Printf(DebugColor, "Debug")
	fmt.Println("")

	InfoLogger.Printf("\n\n#######################################################################################\n|\tNew Instance of algotrading-exchange-manager v0.4.0 (Rel. Date: 03-Apr-2022)    \n#######################################################################################\n")
}

func InitTradeLogger() {
	logFile, err := os.OpenFile("app/zfiles/log/Trades "+time.Now().Format(time.RFC822)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	mwtl := io.MultiWriter(os.Stdout)
	if err != nil {
		log.Println(err)
	} else {
		mwtl = io.MultiWriter(logFile, os.Stdout)
	}
	TradesLogger = log.New(mwtl, "", log.Ltime)

}

func InitLogger() {
	mw := io.MultiWriter(os.Stdout)
	logFile, err := os.OpenFile("app/zfiles/log/algoExchMgr "+time.Now().Format(time.RFC822)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	} else {
		mw = io.MultiWriter(logFile, os.Stdout)
	}

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
