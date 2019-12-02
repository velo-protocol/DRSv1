package main

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"log"
)

func main() {
	callRedeem()
}

func callRedeem() {

	client, err := vclient.NewDefaultTestNetClient("localhost:8080", "SA2LOK42U2OOSXH5JBH2IFS5OWWVD5GILXU4QT3KEA72ZE3B2JTGFGAS")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	result, err := client.RedeemCredit(context.Background(), vtxnbuild.RedeemCredit{
		AssetCode: "vTHB",
		Issuer:    "GAXCJR23OBPRZT2CCH6HJ4HRWFH7WBVM425SK6KE6BBX36FN2G4WAKGP",
		Amount:    "1",
	})
	if err != nil {
		panic(err)
	}
	log.Println(result.HorizonResult.TransactionSuccessToString())
	log.Println("Asset Code To Be Redeemed: ", result.VeloNodeResult.AssetCodeToBeRedeemed)
	log.Println("Asset Issuer To Be Redeemed: ", result.VeloNodeResult.AssetIssuerToBeRedeemed)
	log.Println("Asset Amount To Be Redeemed: ", result.VeloNodeResult.AssetAmountToBeRedeemed)
	log.Println("Collateral Code: ", result.VeloNodeResult.CollateralCode)
	log.Println("Collateral Issuer: ", result.VeloNodeResult.CollateralIssuer)
	log.Println("Collateral Amount: ", result.VeloNodeResult.CollateralAmount)
}
