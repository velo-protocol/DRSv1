package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/testkit/helper"
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

	result, err := client.RedeemCredit(context.Background(), vtxnbuild.RedeemCredit{
		AssetCode: "kDREAM",
		Issuer:    "GAXKPU22AE22NO7FXSW7GTNJJ6FGN5NQLXWTJGNBF4VOKLXVJ3RROXTI",
		Amount:    "1",
	})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
}
