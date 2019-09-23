xdr:
	go run main.go

submit/velotx:
	grpcurl -plaintext \
		-d '{"signedVeloTxXdr":"$(xdr)"}' \
        localhost:8080  grpc.VeloNode/SubmitVeloTx
