FROM golang:alpine

RUN set -ex \
    && apk update \
    && apk add git

WORKDIR /root
COPY vegeta-varload.go attack.csv ./
RUN go get -v -d .
RUN go build vegeta-varload.go

ENTRYPOINT ["/root/vegeta-varload"]
