FROM golang:alpine

RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
    bash \
    curl \
    git \
    nodejs \
    npm \
    openssh

RUN mkdir /git-annex
ENV PATH="${PATH}:/git-annex/git-annex.linux"
RUN apk add --no-cache git openssh curl
RUN curl -Lo /git-annex/git-annex-standalone-amd64.tar.gz https://downloads.kitenet.net/git-annex/linux/current/git-annex-standalone-amd64.tar.gz
RUN cd /git-annex && tar -xzf git-annex-standalone-amd64.tar.gz && rm git-annex-standalone-amd64.tar.gz
RUN apk del --no-cache curl
RUN ln -s /git-annex/git-annex.linux/git-annex-shell /bin/git-annex-shell

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
