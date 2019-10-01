xdr/whitelist:
	go run ./whitelist_op/main.go

xdr/setup-credit:
	go run ./setup_credit_op/main.go

xdr/price-update:
	go run ./price_update_op/main.go

xdr/sign:
	go run ./sign_transaction/main.go


submit/velotx:
	grpcurl -plaintext \
		-d '{"signedVeloTxXdr":"$(xdr)"}' \
        localhost:8080  grpc.VeloNode/SubmitVeloTx
