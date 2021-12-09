Algo Trading Containers [algotrading-ticker-service - algotrading-analysis-service - algotrading-trade-manager]

This container is first part of group of 3 containers to perform algo trading. This container currently supports Zerodha Kite API.

It performs;

    Login
    Token ID generation. Supports NSE (Eq and FUT) and MCX FUT
    Subscribe to Ticks as specified by Token file.
    Save each tick into TimescableDb (To be created as another docker container)

To be done

    Save 1-min Candle based on tick on separate table

----------------- END --------------------------

Source code: https://github.com/parag-b/goTicker
Visit github project page for documentation support


## ENV_Settings.env

TFA_AUTH = ""
USER_ID =""
PASSWORD = ""

API_KEY = ""
API_SECRET = ""
REQUEST_TOKEN_URL ="https://kite.zerodha.com/connect/login?v=3&api_key="

DATABASE_URL = "postgres://username:password@localhost:5432/database_name"
or
DATABASE_URL = "postgres://username:password@abc.com:5432/database_name"

## ENV_accesstoken.env

accessToken=""

# Compile GO

go build
with VC

# Build docker container (and run)
DOCKER_BUILDKIT=1 docker build -t paragba/algotrading-ticker-service:v0.1 .
DOCKER_BUILDKIT=1 docker build -t paragba/algotrading-ticker-service:latest .
docker run --rm -it paragba/algotrading-ticker-service

docker save goticker:latest -o goTickerv0.xx.tar

# Publish to Docker
docker push paragba/algotrading-ticker-service:latest

# TODO Master list
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
