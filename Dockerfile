# BUILDER IMAGE
FROM golang:alpine AS binbuilder

# Build package deps
RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
    bash \
    curl \
    git \
    openssh

# Download git-annex to builder and extract
RUN mkdir /git-annex
RUN curl -Lo /git-annex/git-annex-standalone-amd64.tar.gz https://downloads.kitenet.net/git-annex/linux/current/git-annex-standalone-amd64.tar.gz
RUN cd /git-annex && tar -xzf git-annex-standalone-amd64.tar.gz && rm git-annex-standalone-amd64.tar.gz

ENV GOPROXY https://proxy.golang.org
RUN go version
COPY ./go.mod ./go.sum /gin-valid/
WORKDIR /gin-valid

# download deps before bringing in the sources
RUN go mod download
COPY ./cmd /gin-valid/cmd/
COPY ./internal /gin-valid/internal/
RUN go build ./cmd/ginvalid

### ============================ ###

# RUNNER IMAGE
FROM alpine:latest

# Runtime deps
RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
        && apk --no-cache --no-progress add \
        bash \
        git \
        nodejs \
        npm \
        openssh

# Copy git-annex from builder image
COPY --from=binbuilder /git-annex /git-annex
ENV PATH="${PATH}:/git-annex/git-annex.linux"

# Install the BIDS validator
RUN npm install -g bids-validator

RUN mkdir -p /gin-valid/results/
RUN mkdir -p /gin-valid/tmp/
RUN mkdir -p /gin-valid/config
RUN mkdir -p /gin-valid/tokens/by-sessionid
RUN mkdir -p /gin-valid/tokens/by-repo

ENV GINVALIDHOME /gin-valid/
WORKDIR /gin-valid
ENV GIN_CONFIG_DIR /gin-valid/config/client

# Copy binary and resources into runner image
COPY --from=binbuilder /gin-valid/ginvalid /

ENTRYPOINT /ginvalid --config=/gin-valid/config/cfg.json
EXPOSE 3033
