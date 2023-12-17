FROM golang:1.21

RUN apt-get install tzdata -y

WORKDIR /go/src/app
COPY . .

RUN go get ./...
RUN go build .

ENV CONFIG_FILE=/data/config.json
ENV TZ=America/New_York
VOLUME [ "/data" ]

CMD ["/go/src/app/main"]
