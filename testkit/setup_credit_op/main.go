package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/testkit/helper"
	"log"
)

func main() {
	callSetupCredit()
}

func callSetupCredit() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", helper.SecretKeyTrustedPartner)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.SetupCredit(context.Background(), vtxnbuild.SetupCredit{
		PeggedValue:    "1.0",
		PeggedCurrency: "USD",
		AssetCode:      "vUSD",
	})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
}
