package main

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/testkit/helper"
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
