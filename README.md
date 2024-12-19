# solaredge
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/solaredge?color=green&label=Release&style=plastic)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/solaredge?style=plastic)
![Build](https://github.com/clambin/solaredge/workflows/Test/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/solaredge)
![GitHub](https://img.shields.io/github/license/clambin/solaredge?style=plastic)
[![GoDoc](https://pkg.go.dev/badge/github.com/clambin/solaredge?utm_source=godoc)](http://pkg.go.dev/github.com/clambin/solaredge)

## ⚠️ Breaking Changes in v2.0.0
v2 is a re-implementation of the original client. The main aim of the rewrite is to improve testability of clients.
It is, however, still an implementation of the SolarEdge v1 API.

## Overview
This package provides a client library for the SolarEdge Cloud-Based Monitoring Platform. The API gives access
to data saved in the monitoring servers for your installed SolarEdge equipment and its performance (i.e. generated power & energy).

The implementation is based on SolarEdge's official [API documentation].

The current version of this library implements the following sections of the API:

- Site Data API
- Site Equipment API
- API Versions

Access to SolarEdge data is determined by the user's API Key & installation. If your situation gives you access
to the Accounts List, Meters or Sensors API, feel free to get in touch to get these implemented in this library.

[API documentation]: https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

## Limitations

The following sections of the API have not yet been implemented:

- Account List API
- Meters API
- Sensors API

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
