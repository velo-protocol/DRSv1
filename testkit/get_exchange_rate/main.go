package main

import (
	"context"
	cenGrpc "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/client"
	"log"
)

func main() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", "SAXIAVMNDS4VOM3W7J55T36V3DLY2PKDFHFCQWJTXE5JLUCBONJGUUBM")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	exchangeRate, err := client.GetExchangeRate(context.Background(), &cenGrpc.GetExchangeRateRequest{
		AssetCode: "kDREAM",
		Issuer:    "GAXKPU22AE22NO7FXSW7GTNJJ6FGN5NQLXWTJGNBF4VOKLXVJ3RROXTI",
	})
	if err != nil {
		panic(err)
	}
	log.Println("Asset Code: ", exchangeRate.AssetCode)
	log.Println("Asset Issuer: ", exchangeRate.Issuer)
	log.Println("Redeemable Collateral: ", exchangeRate.RedeemableCollateral)
	log.Println("Redeemable Price Per Unit: ", exchangeRate.RedeemablePricePerUnit)
}
