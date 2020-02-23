FROM golang:alpine

RUN set -ex \
    && apk update \
    && apk add git

WORKDIR /root
COPY vegeta_varload.go attack.csv ./
RUN go get -v -d .
RUN go build vegeta_varload.go

ENTRYPOINT ["/root/vegeta_varload"]
