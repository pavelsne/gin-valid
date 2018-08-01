FROM ubuntu:18.04

RUN apt update
RUN apt install -y \
    wget \
    git \
    git-annex \
    golang \
    nodejs \
    npm

RUN npm install -g bids-validator

RUN mkdir /go
ENV GOPATH /go

RUN mkdir /gin-cli
WORKDIR /gin-cli
RUN wget https://web.gin.g-node.org/G-Node/gin-cli-releases/raw/master/gin-cli-latest-linux-amd64.tar.gz
RUN tar -xf gin-cli-latest-linux-amd64.tar.gz
RUN ln -s /gin-cli/gin /bin/gin

RUN mkdir -p /go/src/github.com/G-Node
ADD . /go/src/github.com/G-Node/gin-valid
RUN cd /go/src/github.com/G-Node/gin-valid
WORKDIR /go/src/github.com/G-Node/gin-valid

RUN go get ./...
RUN go build

RUN mkdir -p /results/gin-valid
RUN mkdir -p /temp
RUN mkdir -p /config

VOLUME ["/results"]
VOLUME ["/temp"]
VOLUME ["/config"]

ENV GINVALIDHOME /results
ENV GINVALIDTEMP /temp

ENTRYPOINT ./gin-valid --config=/config/cfg.json
EXPOSE 3033
