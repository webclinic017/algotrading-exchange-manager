
FROM golang:1.16-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY app/ ./app
COPY *.go ./

RUN go build -o /zerodha_kite_ticker

CMD [ "/zerodha_kite_ticker" ]