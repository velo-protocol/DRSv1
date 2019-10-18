package main

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/testkit/helper"
	"log"
)

func main() {
	callRedeem()
}

func callRedeem() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", helper.SecretKeyRedeemer)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	txResult, err := client.RedeemCredit(context.Background(), vtxnbuild.RedeemCredit{
		AssetCode: "kDREAM",
		Issuer:    "GAXKPU22AE22NO7FXSW7GTNJJ6FGN5NQLXWTJGNBF4VOKLXVJ3RROXTI",
		Amount:    "1",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())
}
