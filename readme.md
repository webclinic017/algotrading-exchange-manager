# Algo Trading  (Ticker #1/3)
## [ algotrading-ticker-service | *algotrading-analysis-service* | *algotrading-trade-manager* ]

This container is first part of 3 micro-services to perform algo trading. This container currently supports **Zerodha Kite API for NSE EQ, NSE FUTs and MCX FUTs**.

### Features
- Auto Login
- Ticker service registers websocket connection at 9am and closes 4pm on weekdays
- Subscribe to instruments as specified by Token file.
- Daily token ID generation/fetch for NSE (Eq and FUT) and MCX FUT
- Ticks are saved in Timescable DB (Use the docker compose provided to spawn)
- *-FUT, *-IDX, *-EQ - Attributes to identify Futures, Index and Equity instruments respectively in DB
    

### To Do
- Create 1-min Candle based on tick in separate DB table
- MCX Silver instrument token generation

# How to use
1. Use the docker-compose file
2. Setup the env variable Zerodha Kite and Database settings
3. Ensure env variable {PRODUCTION: 'true'} & Timezone is set correctly

# Instrument Symbols/Tokens
The symbols to be registered for ticks are stored in trackSymbols.txt
Default token/Instruments file is stored at app/config/trackSymbols.txt
For Futures, as the contract names changes, the name is generated based on todays date.
Post that the instrument token in read from Instruments file downloaded from Zerodha API.
This tokens are used to register the ticks.


# Docker Compose
    version: '3.4'

    services:
    goticker:
        image: paragba/algotrading-ticker-service:latest
        container_name: goTicker
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
        - ./dockerTest/config:/app/app/config
        - ./dockerTest/log:/app/app/log
        
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

# ToDo List
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
