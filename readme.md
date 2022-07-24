# Algo Trading
## [ algotrading-exchange-manager | *algotrading-analysis-service* | *algotrading-trade-manager* ]

This container is first part of 3 micro-services to perform algo trading. This container currently supports **Zerodha Kite API for NSE EQ, NSE FUTs and MCX FUTs**.
### Documentation
https://parag-b.github.io/algotrading-exchange-manager/

### Features
- Auto Login
- Ticker service registers websocket connection at 9am and closes 4pm on weekdays
- Subscribe to instruments as specified by Token file.
- Ticks are saved in Timescable DB (Use the docker compose provided to spawn)
- *-FUT : Attributes to identify Futures instruments in DB
- Compression Policy applied - Data older than 30 days will be auto compressed. Saves up to 90% of disk space.
    

# How to use
1. Use the docker-compose file
2. Setup the env variable Zerodha Kite and Database settings
3. Ensure env variable {PRODUCTION: 'true'} & Timezone is set correctly

# Docker Compose
    version: '3.4'

    services:
    goticker:
        image: paragba/algotrading-exchange-manager:latest
        container_name: algotrading-exchange-manager
        restart: unless-stopped
        environment:
            TZ: 'Asia/Kolkata'
            ZERODHA_USER_ID         : ""
            ZERODHA_PASSWORD        : ""
            ZERODHA_API_KEY         : ""
            ZERODHA_API_SECRET      : ""
            ZERODHA_PIN             : "" # Set ZERODHA_TOTP_SECRET_KEY = `NOT-USED` , to use ZERODHA_PIN for twauth
            ZERODHA_TOTP_SECRET_KEY : "NOT-USED"
            ZERODHA_REQ_TOKEN_URL   : "https://kite.zerodha.com/connect/login?v=3&api_key="
            APP_LIVE_TRADING_MODE   : "TRUE"
            TIMESCALEDB_USERNAME    : ""
            TIMESCALEDB_PASSWORD    : ""
            TIMESCALEDB_ADDRESS     : ""
            TIMESCALEDB_PORT        : ""
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


Source code: https://github.com/parag-b/algotrading-exchange-manager
`Visit github project page for documentation support `


---
## Development
**Compilation** - `go build main.go`

**Create Docker Image** - `DOCKER_BUILDKIT=1 docker build -t paragba/algotrading-exchange-manager:latest .`

**Run Docker** `docker run --rm -it paragba/algotrading-exchange-manager`

**Enter Docker shell** `docker exec -it algo-ex-mgr sh`

**Update libs** `go get -u all && go mod tidy`
**Update specific lib** `go get -u gopkg.in/yaml.v2`

**Generate Coverage Report** 
cd app/trademgr
go tool cover -html=coverage.out
go tool cover -html=coverage.out

---
# Version History


## Version : In-Development
- [ ] API Signal structure modified. OrderBook table updated
- [ ] Candle table - backtesting
- [ ] Post candle cleanup
- [ ] Evaluate movement of all *-FUT to ticks_nsefut table?
- [ ] fast data (buffer) analysis. [tick based strategies]
- [ ] DB buffer size optimisation
- [ ] Loggin file creation - clean ups and new daily file logic
- [ ] Order placement - local

## Version : v0.5.1
- [x] Candles - 1 min candle timescaledb.view (3 day period) created. scheduled everyday @ 5pm
- [x] Strategies - read live data
- [x] Order placement - live testing
- [x] Order tracking
- [x] Order reconcillation
- [x] Order completion/exit
- [x] log files each day for trade and ticker

- Hypertable periods and Compression policy updated for ticks_nsefuts
- DB table names updated
- WdgMgr logic removed
- TradeMgr - Basic core implementation started

## Version : v0.4.3
- [x] all *-FUT are moved to ticks_nsefut table
- [x] Scheduler/CronJob - 8:00 - instruments, 16:05 - candles
- [x] Docker - multistage with distroless image

## Version : v0.4.2

#### Logic's
- Two sepreate tables to store ticks
- nifty-futs will be preserved as ticks data
- rest ticks will be converted to 1-min candle and delete org data
- Trading symbols are read from table (file based support removed)
- Derivative Names - Name calculation replaced with data fetch from instruments file
- Logs are created for each day. Two seperate files for general logs and trade specific logs
- Docker - Command and internal folder structure updated.

 #### Ticker 
- [x] TOTP support
- [x] Two channels created for data transfer
- ticks_nsefut > 'NIFTY-FUT' exlcusive table for storage 
- ticks_stk > table for rest info

#### Trading signals
- [x] Read from tracking_symbols table
- [x] Instruments file - loaded via API everyday, used for fetching tokens

#### Place order
- [x] Margin calculation
- Symbol names - now fetch from instruments table, removed calcualtion based logic 

#### Env variables
- Structure defined in appdata
- [x] Loop parsing to check all data
- Data loaded in struct for application to process
- [ ] Only 'kiteaccessToken' is loaded from os.env()

## Version : ~ v0.3.1
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
