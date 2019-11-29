// +build !unit

package vclient

import (
	"context"
	"github.com/stellar/go/clients/horizonclient"
	cenGrpc "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"log"
)

func ExampleClient_MintCredit() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	whitelistResult, err := client.MintCredit(context.Background(), vtxnbuild.MintCredit{
		AssetCodeToBeIssued: "<ASSET CODE FOR ISSUED>",        // Ex: vTHB
		CollateralAssetCode: "<COLLATERAL ASSET CODE>",        // Ex: VELO (Now Supported Only VELO Token)
		CollateralAmount:    "<COLLATERAL AMOUNT FOR ISSUED>", // Ex: 100
	})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println("Horizon Transaction Hash: ", whitelistResult.HorizonResult.Hash)
	log.Println("Asset Code to be issued: ", whitelistResult.VeloNodeResult.AssetCodeToBeIssued)
	log.Println("Asset Issuer to be issued: ", whitelistResult.VeloNodeResult.AssetIssuerToBeIssued)
	log.Println("Asset Distributor to be issued: ", whitelistResult.VeloNodeResult.AssetDistributorToBeIssued)
	log.Println("Asset Code to be issued: ", whitelistResult.VeloNodeResult.AssetCodeToBeIssued)
	log.Println("Collateral Asset Code: ", whitelistResult.VeloNodeResult.CollateralAssetCode)
	log.Println("Collateral Amount: ", whitelistResult.VeloNodeResult.CollateralAmount)
	// Output:
	// Horizon Transaction Hash: 8d6befa5ddf8845fb75748c81ba360de9728d5e253c8f864b3aef0c1748fa9f5
	// Asset Code to be issued: vTHB
	// Asset Issuer to be issued: GBIO46CY6F2NDON6ZGXHF6WU3P2O7TDW7QREQYAS3OEKFC2RFUYXDSBG
	// Asset Distributor to be issued: GBIO46CY6F2NDON6ZGXHF6WU3P2O7TDW7QREQYAS3OEKFC2RFUYXDSBG
	// Collateral Asset Code: VELO
	// Collateral Amount: 100.0000000

}

func ExampleClient_PriceUpdate() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	whitelistResult, err := client.PriceUpdate(context.Background(), vtxnbuild.PriceUpdate{
		Asset:                       "<COLLATERAL ASSET CODE>",            // Ex: VELO (Now Supported Only VELO Token)
		Currency:                    "<FIAT CURRENCY>",                    // Ex: THB (Now supported Only THB, USD and SGD)
		PriceInCurrencyPerAssetUnit: "<PRICE IN CURRENCY PER ASSET UNIT>", // Ex: 1
	})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println("Horizon Transaction Hash: ", whitelistResult.HorizonResult.Hash)
	log.Println("Currency: ", whitelistResult.VeloNodeResult.Currency)
	log.Println("Collateral Code: ", whitelistResult.VeloNodeResult.CollateralCode)
	log.Println("Price in Currency per Asset Unit: ", whitelistResult.VeloNodeResult.PriceInCurrencyPerAssetUnit)
	// Output:
	// Horizon Transaction Hash: 38cd8ff3e7961efc183a912a8e65d008625b0a88e0482aa14423b21834a8d4ab
	// Currency: THB
	// Collateral Code: VELO
	// Price in Currency per Asset Unit: 1.0000000

}

func ExampleClient_RebalanceReserve() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	whitelistResult, err := client.RebalanceReserve(context.Background(), vtxnbuild.RebalanceReserve{})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println("Horizon Transaction Hash: ", whitelistResult.HorizonResult.Hash)
	//  Output:
	//  Horizon Transaction Hash: 679a30785699e0fa05392fd6f2bae289a381fb2d4ea57571540c1fb87de47515
}

func ExampleClient_RedeemCredit() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	whitelistResult, err := client.RedeemCredit(context.Background(), vtxnbuild.RedeemCredit{
		AssetCode: "<STABLE CREDIT ASSET CODE>",      // Ex: vTHB
		Issuer:    "<STABLE CREDIT ISSUER ADDRESS>",  // Ex: GBIO46CY6F2NDON6ZGXHF6WU3P2O7TDW7QREQYAS3OEKFC2RFUYXDSBG
		Amount:    "<STABLE CREDIT REDEEMED AMOUNT>", // Ex: 50 (Amount of vTHB)
	})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println("Horizon Transaction Hash: ", whitelistResult.HorizonResult.TransactionSuccessToString())
	log.Println("Asset Code to be Redeemed: ", whitelistResult.VeloNodeResult.AssetCodeToBeRedeemed)
	log.Println("Asset Issuer to be Redeemed: ", whitelistResult.VeloNodeResult.AssetIssuerToBeRedeemed)
	log.Println("Asset Amount to be Redeemed: ", whitelistResult.VeloNodeResult.AssetAmountToBeRedeemed)
	log.Println("Collateral Code: ", whitelistResult.VeloNodeResult.CollateralCode)
	log.Println("Collateral Amount: ", whitelistResult.VeloNodeResult.CollateralAmount)
	log.Println("Collateral Issuer: ", whitelistResult.VeloNodeResult.CollateralIssuer)
	// Output:
	// Horizon Transaction Hash: a3b3cd29a04971635a54e59ef62adbbe4400a7465bb989c27b3f5ec332215f1d
	// Asset Code to be Redeemed:  vTHB
	// Asset Issuer to be Redeemed: GBIO46CY6F2NDON6ZGXHF6WU3P2O7TDW7QREQYAS3OEKFC2RFUYXDSBG
	// Asset Amount to be Redeemed: 50.0000000
	// Collateral Code: VELO
	// Collateral Amount: 50.0000000
	// Collateral Issuer: GCNMY2YGZZNUDMHB3EA36FYWW63ZRAWJX5RQZTZXDLRWCK73H77F264J

}

func ExampleClient_SetupCredit() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	setupCreditResult, err := client.SetupCredit(context.Background(), vtxnbuild.SetupCredit{
		PeggedValue:    "<PEGGED VALUE OF STABLE CREDIT>",    // Ex: 1
		PeggedCurrency: "<PEGGED CURRENCY OF STABLE CREDIT>", // Ex: THB
		AssetCode:      "<ASSET CODE OF STABLE CREDIT>",      // Ex: vTHB
	})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println("Horizon Transaction Hash: ", setupCreditResult.HorizonResult.Hash)
	log.Println("Pegged Currency: ", setupCreditResult.VeloNodeResult.PeggedCurrency)
	log.Println("Pegged Value: ", setupCreditResult.VeloNodeResult.PeggedValue)
	log.Println("Asset Code: ", setupCreditResult.VeloNodeResult.AssetCode)
	log.Println("Asset Issuer: ", setupCreditResult.VeloNodeResult.AssetIssuer)
	log.Println("Asset Distributor: ", setupCreditResult.VeloNodeResult.AssetDistributor)

	// Output:
	// Horizon Transaction Hash: a3b3cd29a04971635a54e59ef62adbbe4400a7465bb989c27b3f5ec332215f1d
	// Pegged Currency: THB
	// Pegged Value: 1.0000000
	// Asset Code: vTHB
	// Asset Issuer: GCAXOL5TWF252D32IAWU54WKI2YEUGUYJHBOAKQ4LL6DLFYALNHDJ2E5
	// Asset Distributor: GBP25INM42KH5I4O7YRGGJSYNTF5X6WLVSBS67KSI47YOKALK4LKUGAT

}

func ExampleClient_Whitelist() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	whitelistResult, err := client.Whitelist(context.Background(), vtxnbuild.Whitelist{
		Address:  "<STELLAR ACCOUNT ADDRESS>", // Ex: GDRQ7S53ZZYONK64SG7ABWULW2DYIH4YJAAJEFSQFJDREHJSLDRQIL7I
		Role:     "<VELO ROLE>",               // Ex: PRICE_FEEDER  (Now supported Only TRUSTED_PARTNER, PRICE_FEEDER and REGULATOR)
		Currency: "<FEED CURRENCY>",           // Ex: THB (Now support only THB, USD and SGD), this field is only use for PRICE_FEEDER role
	})
	if err != nil {
		if herr, ok := err.(*horizonclient.Error); ok {
			log.Println(herr.Problem.Detail)
		}
		return
	}

	log.Println("Horizon Transaction Hash: ", whitelistResult.HorizonResult.Hash)
	log.Println("Whitelist Address: ", whitelistResult.VeloNodeResult.Address)
	log.Println("Role: ", whitelistResult.VeloNodeResult.Role)
	log.Println("Feed Currency: ", whitelistResult.VeloNodeResult.Currency)

	// Output:
	// Horizon Transaction Hash: cdb7bb34efaea19fc7ef87be5d77da86923b4fef27023221c336e9bc75bde4ad
	// Whitelist Address: GDRQ7S53ZZYONK64SG7ABWULW2DYIH4YJAAJEFSQFJDREHJSLDRQIL7I
	// Role: PRICE_FEEDER
	// Feed Currency: THB

}

func ExampleClient_GetExchangeRate() {
	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	replyCollateralHealthCheck, err := client.GetExchangeRate(context.Background(), &cenGrpc.GetExchangeRateRequest{
		AssetCode: "<STABLE CREDIT ASSET CODE>",     // Ex: vTHB
		Issuer:    "<STABLE CREDIT ISSUER ADDRESS>", // Ex: GAXKPU22AE22NO7FXSW7GTNJJ6FGN5NQLXWTJGNBF4VOKLXVJ3RROXTI
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Asset Code: ", replyCollateralHealthCheck.AssetCode)
	log.Println("Asset Issuer: ", replyCollateralHealthCheck.Issuer)
	log.Println("RequiredAmount: ", replyCollateralHealthCheck.RedeemableCollateral)
	log.Println("PoolAmount: ", replyCollateralHealthCheck.RedeemablePricePerUnit)

	// Output:
	// Asset Code: vTHB
	// Asset Issuer: GAXKPU22AE22NO7FXSW7GTNJJ6FGN5NQLXWTJGNBF4VOKLXVJ3RROXTI
	// RequiredAmount: 2000.0000000
	// PoolAmount: 1500.0000000

}

func ExampleClient_GetCollateralHealthCheck() {

	client, err := NewDefaultTestNetClient("testnet-drsv1-0.velo.org", clientSecretKey)
	if err != nil {
		log.Println(err)
		return
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
	log.Println("RequiredAmount: ", replyCollateralHealthCheck.RequiredAmount)
	log.Println("PoolAmount: ", replyCollateralHealthCheck.PoolAmount)

	// Output:
	// Asset Code: VELO
	// Asset Issuer: GCNMY2YGZZNUDMHB3EA36FYWW63ZRAWJX5RQZTZXDLRWCK73H77F264J
	// RequiredAmount: 2000.0000000
	// PoolAmount: 1500.0000000

}
