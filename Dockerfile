FROM ubuntu:18.04

RUN apt update
RUN apt install -y \
    git \
    git-annex \
    golang \
    nodejs \
    npm

RUN npm install -g bids-validator

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

ENTRYPOINT gin-valid --config=/config/cfg.json
EXPOSE 3033

COPY ./gin-valid /usr/local/bin/.
