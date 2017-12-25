#### docker build -t golang:1.8.5-ca-certificate .

# FROM golang:1.8.5
#
# RUN apt-get update
# RUN apt-get install -y ca-certificates

#### docker build -t golang:1.8.5-apiserver-boot .

# FROM golang:1.8.5-ca-certificate
#
# RUN wget https://github.com/kubernetes-incubator/apiserver-builder/releases/download/v0.1-alpha.25/apiserver-builder-v0.1-alpha.25-linux-amd64.tar.gz \
#    && tar xvf apiserver-builder-v0.1-alpha.25-linux-amd64.tar.gz  -C /usr/local/bin --strip-components=1 \
#    && rm apiserver-builder-v0.1-alpha.25-linux-amd64.tar.gz

# # For develop branch
# FROM registry.new.dataos.io/datafoundry/golang:1.8.5-apiserver-boot
# 
# COPY . /go/src/github.com/asiainfoldp/apiserver-servicebroker
# WORKDIR /go/src/github.com/asiainfoldp/apiserver-servicebroker
# 
# RUN cp -rf 0-non-gen/pkg/* pkg \
#     && cp -rf 0-non-gen/vendor/* vendor \
#     && apiserver-boot build generated \
#     && apiserver-boot build executables --generate=false \
#     && mv -f bin/* .

# For master branch
FROM registry.new.dataos.io/datafoundry/golang:1.8.5-ca-certificate

COPY . /go/src/github.com/asiainfoldp/apiserver-servicebroker
WORKDIR /go/src/github.com/asiainfoldp/apiserver-servicebroker

RUN go build ./cmd/apiserver \
    && go build ./cmd/controller-manager


