FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y wget
RUN wget -O- http://neuro.debian.net/lists/trusty.de-md.full | tee /etc/apt/sources.list.d/neurodebian.sources.list
RUN apt-get install gnupg -y

RUN apt-key adv --recv-keys --keyserver hkp://pgp.mit.edu:80 0xA5D32F012649A5A9
RUN apt-get update
RUN apt-get install -y \
    git \
    git-annex-standalone\
    golang
RUN apt install -y nodejs
RUN apt install -y npm

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
