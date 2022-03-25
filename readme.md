# Algo Trading
## [ algotrading-exchange-manager | *algotrading-analysis-service* | *algotrading-trade-manager* ]

This container is first part of 3 micro-services to perform algo trading. This container currently supports **Zerodha Kite API for NSE EQ, NSE FUTs and MCX FUTs**.
### Documentation
https://parag-b.github.io/algotrading-exchange-manager/

### Features
- Auto Login
- Ticker service registers websocket connection at 9am and closes 4pm on weekdays
- Subscribe to instruments as specified by Token file.
- Daily token ID generation/fetch for NSE (Eq and FUT) and MCX FUT
- Ticks are saved in Timescable DB (Use the docker compose provided to spawn)
- *-FUT, *-MCXFUT, *-IDX - Attributes to identify Futures, Index and Equity instruments respectively in DB
- 1/3/5/10/15-min Candle tables built-in for faster processing (using Continuous Aggregation of Timescale DB)
- Compression Policy applied - Data before 7 days will be auto compressed. Saves up to 90% of disk space.
    

### To Do

- MCX Silver instrument token generation

# How to use
1. Use the docker-compose file
2. Setup the env variable Zerodha Kite and Database settings
3. Ensure env variable {PRODUCTION: 'true'} & Timezone is set correctly

# Instrument Symbols/Tokens
The symbols to be registered for ticks are stored in trackSymbols.txt
Default token/Instruments file is stored at app/zfiles/config/trackSymbols.txt

For Futures, as the contract names changes, the name is generated based on today's date.
And the current ID is fetched from Instruments file downloaded from Zerodha API.-
This tokens are used to register for the ticks.


# Docker Compose
    version: '3.4'

    services:
    goticker:
        image: paragba/algotrading-exchange-manager:latest
        container_name: algotrading-exchange-manager
        restart: unless-stopped
        environment:
        TZ: 'Asia/Kolkata'
        PRODUCTION: 'true'
        USER_ID: ""
        PASSWORD: ""
        TFA_AUTH: ""  # PIN
        API_KEY: ""
        API_SECRET: ""
        DATABASE_URL: "postgres://postgres:pgpwdChangeMe@mysite.com:5432/stockdb"
        volumes:
        - ./dockerTest/config:/app/app/zfiles/config
        - ./dockerTest/log:/app/app/zfiles/log
        
    # Use the below code to spawn TimescaleDB Container
    timescaledb:
        image: 'timescale/timescaledb:latest-pg12'
        container_name: timescaledb
        restart: unless-stopped
        environment:
        - PUID=1000
        - PGID=1000
        - POSTGRES_PASSWORD=pgpwdChangeMe
        ports:
        - "5432:5432"


Source code: https://github.com/parag-b/goTicker
`Visit github project page for documentation support `



# Development
**Compilation** - `go build main.go`

**Create Docker Image** - `DOCKER_BUILDKIT=1 docker build -t paragba/algotrading-ticker-service:(v0.1/latest) .`

**Run Docker** `docker run --rm -it paragba/algotrading-ticker-service`

**Enter Docker shell** `docker exec -it goTicker sh`

**Update libs** `go get -u all && go mod tidy`
**Update specific lib** `go get -u gopkg.in/yaml.v2`

**Generate Coverage Report** 
cd app/trademgr
go tool cover -html=coverage.out
go tool cover -html=coverage.out

# ToDo List
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
