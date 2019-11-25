package vclient

import (
	"context"
	"github.com/stellar/go/clients/horizonclient"
	cenGrpc "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"log"
)

func ExampleClient_GetCollateralHealthCheck() {

	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = client.Close()
	}()

	replyCollateralHealthCheck, err := client.GetCollateralHealthCheck(context.Background(), &cenGrpc.GetCollateralHealthCheckRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Asset Code: ", replyCollateralHealthCheck.AssetCode)
	log.Println("Asset Issuer: ", replyCollateralHealthCheck.AssetIssuer)
	log.Println("requiredAmount: ", replyCollateralHealthCheck.RequiredAmount)
	log.Println("poolAmount: ", replyCollateralHealthCheck.PoolAmount)
}

func ExampleClient_Whitelist() {
	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		log.Println(err)
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
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println(whitelistResult.VeloNodeResult.Address)
	log.Println(whitelistResult.VeloNodeResult.Role)
	log.Println(whitelistResult.HorizonResult.TransactionSuccessToString())
}

func ExampleClient_SetupCredit() {
	client, err := NewDefaultTestNetClient("localhost:8080", clientSecretKey)
	if err != nil {
		log.Println(err)
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
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println(setupCreditResult.VeloNodeResult.AssetIssuer)
	log.Println(setupCreditResult.VeloNodeResult.AssetCode)
	log.Println(setupCreditResult.HorizonResult.TransactionSuccessToString())
}
