FROM golang:1.18

WORKDIR /go/src/app
COPY . .

RUN go get ./...
RUN go build .

CMD ["/go/src/app/main"]
