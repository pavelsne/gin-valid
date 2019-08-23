# SERVICE BUILDER IMAGE
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

# NIX BUILDER IMAGE
FROM alpine:latest as nixbuilder

# HDF5 is in the 'testing' repository
RUN echo http://dl-2.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories
RUN apk --no-cache --no-progress add \
    git \
    openssh \
    cmake \
    doxygen \
    git \
    build-base \
    boost-dev \
    boost-static \
    cppunit-dev \
    hdf5-dev \
    hdf5-static

RUN git clone https://github.com/G-Node/nix /nix
WORKDIR /nix
RUN git checkout master
RUN mkdir build
WORKDIR build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_EXPORT_COMPILE_COMMANDS=Yes -DBUILD_STATIC=on ..
RUN make all

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

# Install the BIDS validator
RUN npm install -g bids-validator

# Copy git-annex from builder image
COPY --from=binbuilder /git-annex /git-annex
ENV PATH="${PATH}:/git-annex/git-annex.linux"

# Copy nixio-tool from nixbuilder image
COPY --from=nixbuilder /nix/build/nixio-tool /bin

RUN nixio-tool
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
