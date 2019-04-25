[![Build Status](https://travis-ci.org/G-Node/gin-valid.svg?branch=master)](https://travis-ci.org/G-Node/gin-valid)
[![Docker Automated build](https://img.shields.io/docker/automated/gnode/gin-valid.svg)](https://hub.docker.com/r/gnode/gin-valid)

# gin-valid

gin-valid is the G-Node Infrastructure data validation service. It is a microservice server written in go that is meant to be run together with a GIN repository server.

Repositories on a GIN server can trigger validation of data files via this service. Currently validation of the [BIDS](bids.neuroimaging.io) fMRI data format is supported.
