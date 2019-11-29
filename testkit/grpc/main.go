package main

import (
	"fmt"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"github.com/velo-protocol/DRSv1/testkit/helper"
	"log"
)

func main() {
	regulatorKeyPair, _ := vconvert.SecretKeyToKeyPair(helper.SecretKeyFirstRegulator)
	priceFeederKeyPair, _ := vconvert.SecretKeyToKeyPair(helper.SecretKeyPriceFeeder)
	trustedPartnerKeyPair, _ := vconvert.SecretKeyToKeyPair("SB4REBSHGECK7IDDJWTUX3CJSEX5J24O32XWPISQVH5C443L56G6AMQ2")
	redeemerKeyPair, _ := vconvert.SecretKeyToKeyPair(helper.SecretKeyRedeemer)

	// generate whitelist vXDR
	buildB64WhitelistTx(regulatorKeyPair.Address(), helper.PublicKeyPriceFeeder, vxdr.RolePriceFeeder, "THB", regulatorKeyPair)

	// generate price update vXDR
	buildB64PriceUpdateTx(priceFeederKeyPair.Address(), "VELO", "THB", "0.5", priceFeederKeyPair)

	// generate mint vXDR
	buildB64MintTx(trustedPartnerKeyPair.Address(), "vTHB", "VELO", "1", trustedPartnerKeyPair)

	// generate Setup vXDR
	buildB64SetupTx(trustedPartnerKeyPair.Address(), "vTHB", "THB", "1", trustedPartnerKeyPair)

	//  generate Redeem vXDR
	buildB64RedeemTx(redeemerKeyPair.Address(), "vTHB", "GAYOGBGXJVDFRLIAIL4QXSH6ZGRYUDFEZIQTIEGEO6ATSSRG3T4UJ7VE", "10", redeemerKeyPair)

	// generate Rebalance vXDR
	buildB64RebalanceTx(redeemerKeyPair.Address(), redeemerKeyPair)
}

func buildB64WhitelistTx(txSourceAccount, opSourceAccount string, whitelistRole vxdr.Role, currency string, secretKey *keypair.Full) {
	fmt.Println("##### Start Build Whitelist Transaction #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.Whitelist{
			Address:  opSourceAccount,
			Role:     string(whitelistRole),
			Currency: currency,
		},
	}).BuildSignEncode(secretKey)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Whitelist Transaction #####")
}

func buildB64PriceUpdateTx(txSourceAccount, asset, currency, priceInCurrencyPerAssetUnit string, secretKey *keypair.Full) {
	fmt.Println("##### Start Build Price Update Transaction #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.PriceUpdate{
			Asset:                       asset,
			Currency:                    currency,
			PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
		},
	}).BuildSignEncode(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Price Update Transaction #####")

}

func buildB64MintTx(txSourceAccount, assetCodeToBeIssued, collateralAssetCode, collateralAmount string, secretKey *keypair.Full) {
	fmt.Println("##### Start Build Mint Transaction #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.MintCredit{
			AssetCodeToBeIssued: assetCodeToBeIssued,
			CollateralAmount:    collateralAmount,
			CollateralAssetCode: collateralAssetCode,
		},
	}).BuildSignEncode(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Mint Transaction #####")

}

func buildB64SetupTx(txSourceAccount, assetCode, peggedCurrency, peggedValue string, secretKey *keypair.Full) {
	fmt.Println("##### Start Build Setup Transaction #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.SetupCredit{
			PeggedCurrency: peggedCurrency,
			PeggedValue:    peggedValue,
			AssetCode:      assetCode,
		},
	}).BuildSignEncode(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Setup Transaction #####")

}

func buildB64RedeemTx(txSourceAccount, assetCode, issuer, amount string, secretKey *keypair.Full) {
	fmt.Println("##### Start Build Redeem Transaction #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.RedeemCredit{
			AssetCode: assetCode,
			Issuer:    issuer,
			Amount:    amount,
		},
	}).BuildSignEncode(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Redeem Transaction #####")

}

func buildB64RebalanceTx(txSourceAccount string, secretKey *keypair.Full) {
	fmt.Println("##### Start Build Rebalance Transaction #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.RebalanceReserve{},
	}).BuildSignEncode(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Rebalance Transaction #####")

}
