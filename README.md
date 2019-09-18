#Velo Centralized
1. Node - the deployable Velo Node
2. Libraries - xdr encoding/decoding, errors, types

## Project Structure
```
.
|-- grpc
|-- libs
|   |-- crypto
|   |   `-- crypto.go
|   |-- errors
|   |   `-- errors.go
|   |-- types
|   |   `-- types.go
|   `-- xdr
|       `-- xdr.go
|-- node
|   |-- app
|   |   |-- constants
|   |   |   `-- strings.go
|   |   |-- entities
|   |   |   `-- credit.go
|   |   |-- environments
|   |   |   `-- env.go
|   |   |-- extensions
|   |   |   |-- horizon.go
|   |   |   `-- postgres.go
|   |   |-- layers
|   |   |   |-- deliveries
|   |   |   |   `-- grpc
|   |   |   |       |-- check.go
|   |   |   |       |-- init.go
|   |   |   |       `-- watch.go
|   |   |   |-- repositories
|   |   |   |   |-- stellar
|   |   |   |   |   |-- build_mint_tx.go
|   |   |   |   |   |-- build_setup_tx.go
|   |   |   |   |   |-- init.go
|   |   |   |   |   |-- interface.go
|   |   |   |   |   |-- load_account.go
|   |   |   |   |   `-- submit_transaction.go
|   |   |   |   `-- whitelist
|   |   |   `-- usecases
|   |   |       |-- init.go
|   |   |       |-- interface.go
|   |   |       |-- mint.go
|   |   |       |-- setup_account.go
|   |   |       `-- setup_account_test.go
|   |   |-- migrations
|   |   |-- test_helpers
|   |   |   |-- data.go
|   |   |   |-- init.go
|   |   |   `-- stellar.go
|   |   |-- tmp
|   |   |   `-- runner-build
|   |   |-- utils
|   |   |   |-- keypair.go
|   |   |   |-- keypair_test.go
|   |   |   |-- response.go
|   |   |   `-- response_test.go
|   |   `-- main.go
|   |-- resources
|   |   `-- docker
|   |       |-- postgres
|   |       |   `-- initdb.sh
|   |       `-- docker-compose.yaml
|   |-- tmp
|   |-- Dockerfile
|   `-- Makefile
|-- README.md
|-- go.mod
`-- go.sum
```


###Troubleshoot
- If `go mod download` or `go mod tidy` fail, try
```bash
brew install hg
brew install bzr
```

- To consume `cen/node`'s api via command line, it is recommended to use grpCurl.
```bash
grpcurl -plaintext \
    -d '{"signedVeloTxXdr":"AAAA..."}' \
    localhost:8080  grpc.VeloNode/SubmitVeloTx

# if reflection api is not enabled, supply --proto flag
grpcurl -plaintext \
    -d '{"signedVeloTxXdr":"AAAA..."}' \
    --proto ./grpc/velo_node.proto \
    localhost:8080  grpc.VeloNode/SubmitVeloTx
```

