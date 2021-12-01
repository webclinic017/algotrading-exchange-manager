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
DOCKER_BUILDKIT=1 docker build -t goticker .
docker run --rm -it  goticker:latest        

# TODO Master list
- [x] Connect to DB
- [x] Connect to Kite
- [ ] Setup ticker
- [ ] Setup message structure
- [ ] Setup message queue
- [ ] Setup queue consumer
- [ ] Setup queue producer
- [ ] Setup message handler
