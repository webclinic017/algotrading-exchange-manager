#### USE AT OWN RISK. FOR DEMO PURPOSE


# Algo Trading

*Complete algotrading architecture using micro-services*. 

**algotrading-exchange-manager** - Connects to Zerodha API. For **Ticker & Trading**.

**algotrading-analysis-service** - Python based data analysis, implements **Strategies**.

**algotrading-user-portal** - User Interface

### Features
- Supports Zerodha Kite for NSE Equity, Futures and Options - For Ticker & Trading
- Ticker service starts at 9am and closes 4pm on weekdays. Holiday/Spl days not conisidered
- Ticks are saved in Timescale DB
- 1 Min candle table separate created for faster back-testing and reduced RAM consumption
- *-FUT - Attribute to identify Futures instruments
- Complete auto login

# How to use
**Documentation** - https://parag-b.github.io/algotrading-exchange-manager/

**Source code**: https://github.com/parag-b/goTicker


# Docker Compose
    version: '3.4'

    networks:
    gatekeeper_nw:
        external: false
        name: gatekeeper_nw

    services:
    algotrading-exchange-manager:  # !!! This service need to start with delay after DB starts, else issues connecting to DB at init()
        image: paragba/algotrading-exchange-manager:v0.5.4
        container_name: algotrading-exchange-manager
        depends_on:
        algotrading-db:
            condition: service_healthy
        restart: unless-stopped
        networks:
        - gatekeeper_nw    
        environment:
        TZ: 'Asia/Kolkata'
        ZERODHA_LIVE_TEST       : "FALSE"                   # used for UnitTesting with Live trades placd on Zerodha
        APP_LIVE_TRADING_MODE   : "TRUE"                    # not used 
        DB_TBL_TICK_NSEFUT      : "ticks_nsefut"            # DB table names
        DB_TBL_TICK_NSESTK      : "ticks_nsestk"
        DB_TBL_USER_SYMBOLS     : "_symbols"
        DB_TBL_USER_SETTING     : "_setting"
        DB_TBL_USER_STRATEGIES  : "_strategies"
        DB_TBL_ORDER_BOOK       : "_order_book"
        DB_TBL_INSTRUMENTS      : "ticks_instr"
        DB_TBL_CDL_VIEW_STK     : "view_1min_cdl_stk"
        DB_TBL_CDL_VIEW_FUT     : "view_1min_cdl_fut"      
        DB_TEST_PREFIX          : "_test"
        DB_TBL_PREFIX_USER_ID   : "myname"    
        ZERODHA_PASSWORD        : ""
        ZERODHA_USER_ID         : ""
        ZERODHA_API_KEY         : ""
        ZERODHA_API_SECRET      : ""
        ZERODHA_PIN             : ""
        ZERODHA_TOTP_SECRET_KEY : ""                      # Key provided by Zerodha while enabling TOTP. Mandatory for trading API.
        ZERODHA_REQ_TOKEN_URL   : "https://kite.zerodha.com/connect/login?v=3&api_key="
        TIMESCALEDB_USERNAME    : "postgres"
        TIMESCALEDB_PASSWORD    : "pgpwd"
        TIMESCALEDB_ADDRESS     : "algotrading-db"
        TIMESCALEDB_PORT        : "5432"
        ALGO_ANALYSIS_ADDRESS   : "http://algotrading-analysis-service:5000/"
        volumes:
        - /algotrading/algoExchMgr:/app/zfiles/log
        logging:
            driver: "json-file"
            options:
            max-file: "5"   # number of files or file count
            max-size: "10m" # file size    
    
    # Use the below code to spawn TimescaleDB Container

    algotrading-db:
        image: 'timescale/timescaledb:latest-pg14'
        container_name: algotrading-db
        restart: unless-stopped
        networks:
        - gatekeeper_nw         
        environment:
        - TZ='Asia/Kolkata'
        - PUID=1000
        - PGID=1000
        - POSTGRES_PASSWORD=pgpwd
        ports:
        - "5405:5432"
        volumes:
        - /algotrading/algotrading-db:/var/lib/postgresql/data
        healthcheck:
        test: ["CMD-SHELL", "pg_isready -U postgres"]
        interval: 5s
        timeout: 5s
        retries: 5
        logging:
            driver: "json-file"
            options:
            max-file: "5"   # number of files or file count
            max-size: "10m" # file size          
        
    algotrading-analysis-service:
        image: paragba/algotrading-analysis-service:v0.4.5
        container_name: algotrading-analysis-service
        depends_on:
        algotrading-db:
            condition: service_healthy
        restart: unless-stopped
        environment:
        TZ: 'Asia/Kolkata'
        TIMESCALEDB_NAME        : "algotrading"
        TIMESCALEDB_ADDRESS     : "algotrading-db"
        TIMESCALEDB_USER        : "postgres"
        TIMESCALEDB_PASSWORD    : "pgpwd"
        TIMESCALEDB_PORT        : "5432"
        DB_TBL_PREFIX_USER_ID   : "myname"
        DB_TBL_INSTRUMENTS      : "ticks_instr"
        DB_TBL_TICK_NSEFUT      : "ticks_nsefut"
        DB_TBL_TICK_NSESTK      : "ticks_nsestk"
        DB_TBL_USER_SETTING     : "_setting"
        DB_TBL_USER_SYMBOLS     : "_symbols"
        DB_TBL_ORDER_BOOK       : "_order_book"
        DB_TBL_USER_STRATEGIES  : "_strategies"
        DB_TBL_CDL_VIEW_STK     : "view_1min_cdl_stk"
        DB_TBL_CDL_VIEW_FUT     : "view_1min_cdl_fut"
        DB_TEST_POSTFIX         : "_test"
        #ports:
        #- "5001:5000"
        networks:
        - gatekeeper_nw          
        volumes:
        - /algotrading/algo-analysis-service/candle_creator:/usr/local/bin/Data/candle_converter
        logging:
            driver: "json-file"
            options:
            max-file: "5"   # number of files or file count
            max-size: "10m" # file size 

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
- [ ] Candle table - backtesting
- [ ] Post candle cleanup
- [ ] Evaluate movement of all *-FUT to ticks_nsefut table?
- [ ] fast data (buffer) analysis. [tick based strategies]
- [ ] DB buffer size optimisation
- [ ] Loggin file creation - clean ups and new daily/symbol wise file logic

## Version : v0.5.5
- [x] Sell Qty bug fix [Issue #40]
- [x] Trades resumes in 2 states only - "AwaitSignal", "TradeMonitoring"
- [x]  New Order creation in DB, replaced search logic with LASTVAL transaction/commit operation [Issue #36]

## Version : v0.5.4
- [x] 1Min Candles -API created. Copy from view to table. View updated @5pm, copy API invoked at 10pm every working day.

## Version : v0.5.3
- [x] golang v1.18
- [x] API Signal structure modified. OrderBook table updated
- [x] Order placement - setup for live trade

## Version : v0.5.3
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
