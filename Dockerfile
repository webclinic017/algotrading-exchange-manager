
FROM golang:1.16-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY app/ ./app
COPY *.go ./
RUN rm ./app/config/*.env

RUN go build -o /algotrading-exchange-manager

CMD [ "/algotrading-exchange-manager" ]