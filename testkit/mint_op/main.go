package main

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/testkit/helper"
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

	txResult, err := client.MintCredit(context.Background(), vtxnbuild.MintCredit{
		AssetCodeToBeIssued: "kDREAM",
		CollateralAssetCode: "VELO",
		CollateralAmount:    "10",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())
}
