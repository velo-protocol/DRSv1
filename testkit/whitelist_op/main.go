package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/testkit/helper"
	"log"
)

func main() {
	callWhitelist()
}

func callWhitelist() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", helper.SecretKeyFirstRegulator)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.Whitelist(context.Background(), vtxnbuild.Whitelist{
		Address:  "GC5F4E7IKMDFNOL7Z5WDHC42LBLLQL2UFY6KQALO2RRHC5EMJJRECPI3",
		Role:     "PRICE_FEEDER",
		Currency: "USD",
	})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
}
