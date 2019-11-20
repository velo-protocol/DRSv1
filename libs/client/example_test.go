package vclient

import (
	"context"
	"github.com/stellar/go/clients/horizonclient"
	cenGrpc "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
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

	priceUpdateResult, err := client.PriceUpdate(context.Background(), vtxnbuild.PriceUpdate{
		Asset:                       "VELO",
		Currency:                    "USD",
		PriceInCurrencyPerAssetUnit: "0.5",
	})
	if err != nil {
		panic(err)
	}
	log.Println(priceUpdateResult.HorizonResult.TransactionSuccessToString())

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

	replyCollateralHealthCheck, err := client.GetCollateralHealthCheck(context.Background(), &cenGrpc.GetCollateralHealthCheckRequest{})
	if err != nil {
		panic(err)
	}
	log.Println("Asset Code: ", replyCollateralHealthCheck.AssetCode)
	log.Println("Asset Issuer: ", replyCollateralHealthCheck.AssetIssuer)
	log.Println("requiredAmount: ", replyCollateralHealthCheck.RequiredAmount)
	log.Println("poolAmount: ", replyCollateralHealthCheck.PoolAmount)

}

func ExampleClient_Whitelist() {
	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	whitelistResult, err := client.Whitelist(context.Background(), vtxnbuild.Whitelist{
		Address:  "GC5F4E7IKMDFNOL7Z5WDHC42LBLLQL2UFY6KQALO2RRHC5EMJJRECPI3",
		Role:     "PRICE_FEEDER",
		Currency: "USD",
	})
	if err != nil {
		// In case horizon returns error
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}

		// In case Velo Node has no error, you can log for response from Velo
		if whitelistResult.VeloNodeResult != nil {
			log.Println(whitelistResult.VeloNodeResult.Address)
			log.Println(whitelistResult.VeloNodeResult.Role)
		}

		panic(err)
	}

	log.Println(whitelistResult.HorizonResult.TransactionSuccessToString())
}

func ExampleClient_SetupCredit() {
	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()

	setupCreditResult, err := client.SetupCredit(context.Background(), vtxnbuild.SetupCredit{
		PeggedValue:    "1.0",
		PeggedCurrency: "USD",
		AssetCode:      "vUSD",
	})
	if err != nil {
		// In case horizon returns error
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}

		// In case Velo Node has no error, you can log for response from Velo
		if setupCreditResult.VeloNodeResult != nil {
			log.Println(setupCreditResult.VeloNodeResult.AssetIssuer)
			log.Println(setupCreditResult.VeloNodeResult.AssetCode)
		}

		panic(err)
	}

	log.Println(setupCreditResult.HorizonResult.TransactionSuccessToString())
}
