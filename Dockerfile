
FROM golang:1.18-alpine

WORKDIR /
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
RUN go mod download

COPY app/ ./app
RUN rm -f ./app/zfiles/config/*.env
RUN rm -f ./app/zfiles/log/*.log

RUN go build -o /algoexmgr

CMD [ "/algoexmgr" ]
