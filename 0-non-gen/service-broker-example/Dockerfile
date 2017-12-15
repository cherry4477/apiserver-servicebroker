FROM golang:1.8.5

COPY . /go/src/github.com/asiainfoldp/apiserver-servicebroker/0-non-gen/service-broker-example
WORKDIR /go/src/github.com/asiainfoldp/apiserver-servicebroker/0-non-gen/service-broker-example

RUN go build

CMD ["sh", "-c", "./service-broker-example -port=33333"]

