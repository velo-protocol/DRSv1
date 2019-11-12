package main

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/testkit/helper"
	"log"
)

func main() {
	callPriceUpdate()
}

func callPriceUpdate() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", helper.SecretKeyPriceFeeder)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.PriceUpdate(context.Background(), vtxnbuild.PriceUpdate{
		Asset:                       "VELO",
		Currency:                    "THB",
		PriceInCurrencyPerAssetUnit: "0.5",
	})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
}
