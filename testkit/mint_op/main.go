package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/testkit/helper"
	"log"
)

func main() {
	callMint()
}

func callMint() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", helper.SecretKeyTrustedPartner)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.MintCredit(context.Background(), vtxnbuild.MintCredit{
		AssetCodeToBeIssued: "kDREAM",
		CollateralAssetCode: "VELO",
		CollateralAmount:    "10",
	})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
}
