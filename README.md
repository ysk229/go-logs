# go-logs



[![Build Status](https://github.com/ysk229/go-logs/workflows/Go/badge.svg)](https://github.com/ysk229/go-logs/actions)
[![Github License](https://img.shields.io/github/license/ysk229/go-logs.svg?style=flat)](https://github.com/ysk229/go-logs/blob/master/LICENSE)
[![Go Doc](https://godoc.org/github.com/ysk229/go-logs?status.svg)](https://pkg.go.dev/github.com/ysk229/go-logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ysk229/go-logs)](https://goreportcard.com/report/github.com/ysk229/go-logs)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/ysk229/go-logs)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/ysk229/go-logs)
[![Github Latest Release](https://img.shields.io/github/release/ysk229/go-logs.svg?style=flat)](https://github.com/ysk229/go-logs/releases/latest)
[![Github Latest Tag](https://img.shields.io/github/tag/ysk229/go-logs.svg?style=flat)](https://github.com/ysk229/go-logs/tags)
[![Github Stars](https://img.shields.io/github/stars/ysk229/go-logs.svg?style=flat)](https://github.com/ysk229/go-logs/stargazers)

## ⚙️ Installation

Download rabbitmq package by using:
```bash
go get github.com/ysk229/go-logs
```

## 🚀 Quick Start 

### Default options

```go
    l := log.DefaultLogger
    l.Info("info") // 不打印
    l.Warnw("key", "value")
    l.Errorf("%v", "sdfdsfds")
    l.SetLevel("info")
    lg.Println("test")
    l.Info("show info log")
    l.SetLevel("warn")
    l.Warnw("key", "value")
```
 
## Other usage examples

See the [examples](_examples) directory for more ideas.



## Transient Dependencies

- [go.uber.org/zap](https://github.com/uber-go/zap)
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)

## 👏 Contributing

I love help! Contribute by forking the repo and opening pull requests. Please ensure that your code passes the existing tests and linting, and write tests to test your changes if applicable.

All pull requests should be submitted to the `main` branch.
