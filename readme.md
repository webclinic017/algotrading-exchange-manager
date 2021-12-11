> Algo Trading Containers [algotrading-ticker-service - algotrading-analysis-service - algotrading-trade-manager]

This container is first part of group of 3 containers to perform algo trading. This container currently supports **Zerodha Kite API for NSE EQ, NSE FUTs and MCX FUTs**.

It performs;

    * Auto Login
    * Token ID generation for NSE (Eq and FUT) and MCX FUT
    * Subscribe to Ticks as specified by Token file.
    * Save each tick into TimescableDb (To be created as another docker container)
    * Ticker starts tick websocket connection at 9am and closes 4pm on weekdays

To be done

    * Save 1-min Candle based on tick on separate table
    * MCX Silver token generation issue

# How to use
- Use the docker-compose file
- Setup the env variable Zerodha Kite and Database settings
- Ensure env variable {PRODUCTION: 'true'}
- Ensure Timezone is set as per your zone/server

# Settings
- USER_ID =""
- TFA_AUTH = "" // AKA PIN
- PASSWORD = ""
- API_KEY = ""
- API_SECRET = ""
- DATABASE_URL = "postgres://username:password@localhost:5432/database_name"

# Instrument Symbols/Tokens
The symbols to be registered for ticks are stored in trackSymbols.txt
Default token/Instruments file is stored at app/config/trackSymbols.txt
For Futures, as the contract names changes, the name is generated based on todays date.
Post that the instrument token in read from Instruments file downloaded from Zerodha API.
This tokens are used to register the ticks.

Source code: https://github.com/parag-b/goTicker
Visit github project page for documentation support


# Development
**Compilation** - `go build main.go`

**Create Docker Image** - `DOCKER_BUILDKIT=1 docker build -t paragba/algotrading-ticker-service:(v0.1/latest) .`

**Run Docker** `docker run --rm -it paragba/algotrading-ticker-service`

**Enter Docker shell** `docker exec -it goTicker sh`

# ToDo List
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
