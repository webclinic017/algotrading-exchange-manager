> Algo Trading Containers [algotrading-ticker-service - algotrading-analysis-service - algotrading-trade-manager]

This container is first part of group of 3 containers to perform algo trading. This container currently supports **Zerodha Kite API for NSE EQ, NSE FUTs and MCX FUTs**.

It performs;

    * Auto Login
    * Token ID generation. Supports NSE (Eq and FUT) and MCX FUT
    * Subscribe to Ticks as specified by Token file.
    * Save each tick into TimescableDb (To be created as another docker container)

To be done

    * Save 1-min Candle based on tick on separate table

# How to use
- Use the docker-compose file
- On first run, default templates will be created in "config" folder.
- Update the Zerodha Kite and Database setting.
- Ensure Timezone is set as per your zone/server
- Restart the docker.


# Env settings
*Settings are stored in app/config/ENV_Settings.env file*

- TFA_AUTH = ""
- USER_ID =""
- PASSWORD = ""
- API_KEY = ""
- API_SECRET = ""
- REQUEST_TOKEN_URL ="https://kite.zerodha.com/connect/login?v=3&api_key="
- DATABASE_URL = "postgres://username:password@localhost:5432/database_name"
or
- DATABASE_URL = "postgres://username:password@abc.com:5432/database_name"

*Access token from Kite login is stored in app/config/ENV_accessToken.env*

#Token File
*Defulat token/Instruments file is stored at app/config/trackSymbols.txt*

# Development
**Compilation** - `go build main.go`

**Create Docker Image** - `DOCKER_BUILDKIT=1 docker build -t paragba/algotrading-ticker-service:(v0.1/latest) .`

**Run Docker** `docker run --rm -it paragba/algotrading-ticker-service`

**Enter Docker shell** `docker exec -it gotickerTest sh`

# TODO Master list
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
