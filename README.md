
<div align="center">
<a href="https://velo.org"><img alt="Stellar" src="https://raw.githubusercontent.com/velo-protocol/assets/master/images/logo.png" width="368" /></a>
<br/>
<h1>VELO DRS</h1>
</div>

[![GoDoc](https://godoc.org/github.com/velo-protocol/DRSv1?status.svg)](https://godoc.org/github.com/velo-protocol/DRSv1)  [![Go Report Card](https://goreportcard.com/badge/github.com/velo-protocol/DRSv1)](https://goreportcard.com/report/github.com/velo-protocol/DRSv1)

## Package Index
- [Node](node) : the deployable Velo Node
- [Go Velo SDK - vclient](libs/client) : Client for Velo Node (queries and transaction submission)
- [Go Velo SDK - vtxnbuild](libs/txnbuild) : Construct Velo transactions and operations
- [Go Velo CLI]() : Cli for Velo Node
- [Document](https://docs.velo.org/) : Velo documentation
 
### Dependencies
This repository is officially supported on the last two releases of Go, which is currently Go 1.12 and Go 1.13.

It depends on a [number of external dependencies](go.mod), and uses Go [Modules](https://github.com/golang/go/wiki/Modules) to manage them. Running any `go` command will automatically download dependencies required for that operation.

You can choose to checkout this repository into a [GOPATH](https://github.com/golang/go/wiki/GOPATH) or into any directory, but if you are using a GOPATH with Go 1.12 or earlier you must set environment variable `GO111MODULE=on` to enable Modules.

### Troubleshoot
- If `go mod download` or `go mod tidy` fail, try
```bash
brew install hg
brew install bzr
```

