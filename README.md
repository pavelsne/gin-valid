[![Build Status](https://travis-ci.org/G-Node/gin-valid.svg?branch=master)](https://travis-ci.org/G-Node/gin-valid)
[![Docker Automated build](https://img.shields.io/docker/automated/gnode/gin-valid.svg)](https://hub.docker.com/r/gnode/gin-valid)

# gin-valid

gin-valid is the G-Node Infrastructure data validation service. It is a microservice server written in go that is meant to be run together with a GIN repository server.

Repositories on a GIN server can trigger validation of data files via this service. Currently there are two validators supported:
- The [BIDS](https://bids.neuroimaging.io) fMRI data format.
- The [NIX](http://g-node.org/nix) (Neuroscience Information Exchange) format.

## Contributing

For instructions on how to add more validators, see the [adding validators](docs/adding-validators.md) contribution guide.
