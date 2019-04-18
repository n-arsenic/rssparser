FROM golang:1.11

WORKDIR $GOPATH/src/rssparser

COPY . . # не копировать папку vendor

RUN make auto-build

RUN make build

WORKDIR $GOPATH/bin

CMD ["app"]
