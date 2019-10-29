package main

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"log"
)

func main() {
	rebalanceReserve()
}

func rebalanceReserve() {
	client, err := vclient.NewDefaultTestNetClient("localhost:8080", "<YOUR KEY PAIR>")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	txResult, err := client.RebalanceReserve(context.Background(), vtxnbuild.RebalanceReserve{})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())
}
