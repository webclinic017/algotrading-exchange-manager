
FROM golang:1.18-alpine

WORKDIR /
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
RUN go mod download

COPY app/ ./app
<<<<<<< HEAD
RUN rm -f ./app/zfiles/config/*.env
RUN rm -f ./app/zfiles/log/*.log
=======
COPY *.go ./
RUN rm ./app/zfiles/config/userSettings.env
RUN rm ./app/zfiles/log/*.log
>>>>>>> 74c584a6a67a3da340908364baaa08e9062ac9ca

RUN go build -o /algoexmgr

CMD [ "/algoexmgr" ]
