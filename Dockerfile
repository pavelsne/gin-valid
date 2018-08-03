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

# RUN mkdir -p /go/src/github.com/G-Node
# ADD . /go/src/github.com/G-Node/gin-valid
# RUN cd /go/src/github.com/G-Node/gin-valid
# WORKDIR /go/src/github.com/G-Node/gin-valid

# RUN go get -v ./...
# RUN go build

RUN mkdir -p /gin-valid/results/
RUN mkdir -p /gin-valid/tmp/
RUN mkdir -p /gin-valid/config

VOLUME ["/gin-valid/"]

ENV GINVALIDHOME /gin-valid/
ENV GIN_CONFIG_DIR /gin-valid/config/client

COPY ./gin-valid .
RUN mkdir -p /root/.config/g-node/gin/

ENTRYPOINT ./gin-valid --config=/config/cfg.json
EXPOSE 3033
