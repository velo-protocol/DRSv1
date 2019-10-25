package vclient

import (
	"context"
	cenGrpc "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"log"
)

func Example() {
	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	txResult, err := client.Whitelist(context.Background(), vtxnbuild.Whitelist{
		Address:  "GC5F4E7IKMDFNOL7Z5WDHC42LBLLQL2UFY6KQALO2RRHC5EMJJRECPI3",
		Role:     "PRICE_FEEDER",
		Currency: "USD",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())

	txResult, err = client.SetupCredit(context.Background(), vtxnbuild.SetupCredit{
		PeggedValue:    "1.0",
		PeggedCurrency: "USD",
		AssetCode:      "vUSD",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())

	txResult, err = client.PriceUpdate(context.Background(), vtxnbuild.PriceUpdate{
		Asset:                       "VELO",
		Currency:                    "USD",
		PriceInCurrencyPerAssetUnit: "0.5",
	})
	if err != nil {
		panic(err)
	}
	log.Println(txResult.TransactionSuccessToString())

	reply, err := client.GetExchangeRate(context.Background(), &cenGrpc.GetExchangeRateRequest{
		AssetCode: "vUSD",
		Issuer:    "GC5F4E7IKMDFNOL7Z5WDHC42LBLLQL2UFY6KQALO2RRHC5EMJJRECPI3",
	})
	if err != nil {
		panic(err)
	}
	log.Println("Asset Code: ", reply.AssetCode)
	log.Println("Asset Issuer: ", reply.Issuer)
	log.Println("Redeemable Collateral: ", reply.RedeemableCollateral)
	log.Println("Redeemable Price Per Unit: ", reply.RedeemablePricePerUnit)

	replyCollateralHealthCheck, err := client.GetCollateralHealthCheck(context.Background(), &cenGrpc.Empty{})
	if err != nil {
		panic(err)
	}
	log.Println("Asset Code: ", replyCollateralHealthCheck.AssetCode)
	log.Println("Asset Issuer: ", replyCollateralHealthCheck.AssetIssuer)
	log.Println("requiredAmount: ", replyCollateralHealthCheck.RequiredAmount)
	log.Println("poolAmount: ", replyCollateralHealthCheck.PoolAmount)

}
