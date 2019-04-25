FROM ubuntu:19.04

RUN apt update
RUN apt install -y \
    git \
    git-annex \
    golang \
    nodejs \
    npm

RUN npm install -g bids-validator

RUN mkdir -p /gin-valid/results/
RUN mkdir -p /gin-valid/tmp/
RUN mkdir -p /gin-valid/config
RUN mkdir -p /gin-valid/tokens/by-sessionid
RUN mkdir -p /gin-valid/tokens/by-repo

VOLUME ["/gin-valid/"]

ENV GINVALIDHOME /gin-valid/
ENV GIN_CONFIG_DIR /gin-valid/config/client

ENV GOPATH /go

# getting repo version so we have a base snapshot of the upstream gin-valid and
# (more importantly) its dependencies, so that docker doesn't have to download
# everything every time a file changes in the local directory
RUN go get -v github.com/G-Node/gin-valid

COPY . $GOPATH/src/github.com/G-Node/gin-valid
WORKDIR $GOPATH/src/github.com/G-Node/gin-valid

RUN go get -v ./...
RUN go build

ENTRYPOINT ./gin-valid --config=/config/cfg.json
EXPOSE 3033
