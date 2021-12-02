
FROM golang:1.16-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./app ./app

RUN go build -o app/goTicker.go app/goTicker.go

CMD [ "app/goTicker.go" ]