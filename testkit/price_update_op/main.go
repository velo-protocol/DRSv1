package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"log"
)

func main() {
	callPriceUpdate()
}

func callPriceUpdate() {

	client, err := vclient.NewDefaultTestNetClient("dev-velo-cen-node-01.velo.org:8080", "SABZJDPDV3BLYBJD4KZF3ZT4MARRHXFX4VNBMVKZ4NSSXVTXZ3E7H454")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.PriceUpdate(context.Background(), vtxnbuild.PriceUpdate{
		Asset:                       "VELO",
		Currency:                    "USD",
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
