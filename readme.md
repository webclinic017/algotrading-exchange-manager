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
DOCKER_BUILDKIT=1 docker build -t paragba/zerodha_kite_ticker .
docker run --rm -it  goticker:latest        
docker save goticker:latest -o goTickerv0.xx.tar

# Publish to Docker
docker push paragba/zerodha_kite_ticker:goTicker

# TODO Master list
- [x] Connect to DB
- [x] Connect to Kite
- [x] Setup ticker
- [x] Setup message structure
- [x] Setup message queue
- [x] Setup queue consumer
- [x] Setup queue producer
- [x] Setup message handler
