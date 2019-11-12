package main

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/testkit/helper"
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
