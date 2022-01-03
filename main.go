package main

import (
	"goTicker/app/db"
	"goTicker/app/kite"
	"goTicker/app/srv"
	"goTicker/graph"
	"goTicker/graph/generated"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

var (
	envOk, dbOk, kiteOk           bool
	apiKey, accToken              string
	wdg, closeTicker, initTicker  *cron.Cron
	symbolFutStr, symbolMcxFutStr string
	Tokens                        []uint32
)

func main() {

	// timeZone, _ := time.LoadLocation("Asia/Calcutta")

	srv.CheckFiles()
	srv.Init()

	initTickerToken(false) // Check if conections are okay

	// start watchdog to recover from connections issues
	wdg = cron.New()
	wdg.AddFunc("@every 10s", checkConnection)
	wdg.Start()

	// everyday scheduled start At 09:00:00 Mon-Fri
	initTicker = cron.New()
	initTicker.AddFunc("0 0 9 * * 1-5", initalizeKite)
	initTicker.Start()

	// everyday scheduled stop At 16:00:00 Mon-Fri
	closeTicker = cron.New()
	closeTicker.AddFunc("0 0 16 * * 1-5", StopKite)
	closeTicker.Start()

	startGraphQL()

}

func startGraphQL() {

	const defaultPort = "5555"

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadEnv() bool {

	// Load .env file, if not in production

	println("PRODUCTION - ", os.Getenv("PRODUCTION"))
	if os.Getenv("PRODUCTION") != "true" {
		srv.InfoLogger.Println("DEVELOPMENT ENV - Ensure ENV variables are set in ENV_settings.env")
		srv.FileCopyIfMissing("app/templates/ENV_Settings.env", "app/config/ENV_Settings.env")
		_ = godotenv.Load("app/config/ENV_Settings.env")
	} else {
		srv.InfoLogger.Println("PRODUCTION ENV- Ensure ENV variables are set")
	}

	srv.InfoLogger.Println("user ID", os.Getenv("USER_ID"))

	if 0 >= len(os.Getenv("USER_ID")) {
		srv.ErrorLogger.Println("USER_ID not set")
		return false
	}
	if 0 >= len(os.Getenv("TFA_AUTH")) {
		srv.ErrorLogger.Println("TFA_AUTH not set")
		return false
	}
	if 0 >= len(os.Getenv("PASSWORD")) {

		srv.ErrorLogger.Println("PASSWORD not set")
		return false
	}
	if 0 >= len(os.Getenv("API_KEY")) {
		srv.ErrorLogger.Println("API_KEY not set")
		return false
	}
	if 0 >= len(os.Getenv("API_SECRET")) {
		srv.ErrorLogger.Println("API_SECRET not set")
		return false
	}
	if 0 >= len(os.Getenv("DATABASE_URL")) {
		srv.ErrorLogger.Println("DATABASE_URL not set")
		return false
	}
	return true
}

func initTickerToken(check bool) {

	envOk = loadEnv()

	if envOk {

		dbOk = db.DbInit()
		db.StoreSymbolsInDb(symbolFutStr, symbolMcxFutStr)

		kite.Tokens, kite.InsNamesMap, symbolFutStr, symbolMcxFutStr = kite.GetSymbols()

		kiteOk, apiKey, accToken = kite.LoginKite()

		if envOk && dbOk && kiteOk && check {
			// Initate zerodha ticker
			kite.TickerInitialize(apiKey, accToken)
			go db.StoreTickInDb()
			// start watchdog to recover from connections issues
		}
	}
}
func StopKite() {
	kite.CloseTicker()
	db.CloseDBPool()
}

func initalizeKite() {
	kite.CloseTicker()
	srv.InfoLogger.Println("\nInitializing Kite", kiteOk)
	initTickerToken(true)
	srv.InfoLogger.Printf("\n\n\t--------------STATUS---------------\n\t| Environment variables set: %t |\n\t| Kite Login Succesfull: %t     |\n\t| DB Connected: %t              |\n\t-----------------------------------\n\n", envOk, kiteOk, dbOk)
}

func checkConnection() {

	if !kite.KiteConnectionStatus {

		now := time.Now()

		if (now.Hour() >= 9) && (now.Hour() < 16) &&
			(now.Weekday() > 0) && (now.Weekday() < 6) {
			initalizeKite()
		}
	}
}
