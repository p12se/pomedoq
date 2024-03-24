<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/p12se/pomedoq/raw/master/assets/images/logo_white.png">
    <img alt="Pomedoq logo" src="https://github.com/p12se/pomedoq/raw/master/assets/images/logo_dark.png" width="40%">
  </picture>
</div>

<div align="center">

![build](https://github.com/p12se/pomedoq/actions/workflows/ci.yaml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/p12se/pomedoq)
![GitHub License](https://img.shields.io/github/license/p12se/pomedoq)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/p12se/pomedoq)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/p12se/pomedoq)
![GitHub issues](https://img.shields.io/github/issues/p12se/pomedoq)
![GitHub pull requests](https://img.shields.io/github/issues-pr/p12se/pomedoq)
![GitHub contributors](https://img.shields.io/github/contributors/p12se/pomedoq)

</div>

## Pomedoq
Pomedoq is a tool that generates markdown documentation for prometheus metrics. It reads the metrics from a prometheus server and generates a markdown output with the metrics information. The markdown output contains the metric name, description, type, labels, and help text. The generated markdown output can be used as documentation for the prometheus metrics.

## Getting Started
* Install Pomedoq by downloading the binary from the [releases](https://github.com/p12se/pomedoq/releases) page. 

## Install with go install
This method requires [Go](https://go.dev) to be installed on your system.
```bash
go get -u github.com/p12se/pomedoq@latest
```
