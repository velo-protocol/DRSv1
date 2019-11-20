package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/testkit/helper"
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
	log.Println("Collateral Code: ", result.VeloNodeResult.CollateralCode)
	log.Println("Currency: ", result.VeloNodeResult.Currency)
	log.Println("Price In Currency Per Asset Unit: ", result.VeloNodeResult.PriceInCurrencyPerAssetUnit)
}
