
FROM golang:1.18-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY app/ ./app
COPY *.go ./
RUN rm ./app/zfiles/config/userSettings.env
RUN rm ./app/zfiles/log/*.log

RUN go build -o /algoexmgr

CMD [ "/algoexmgr" ]