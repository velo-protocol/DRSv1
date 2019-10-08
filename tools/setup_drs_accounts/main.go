package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	_operations "gitlab.com/velo-labs/cen/tools/setup_drs_accounts/operations"
	"log"
)

var (
	sourceSeedKey          = "<YOUR_SECRET_KEY>"
	veloAssetIssuerAccount = "<VELO_ISSUER_PUBLIC_KEY>"
	veloAsset              = "VELO"

	sourceKP, _ = vconvert.SecretKeyToKeyPair(sourceSeedKey)
	client      = horizonclient.DefaultTestNetClient
)

func main() {
	generateDRSAndFriends()
}

func loadAccount(publicKey string) *horizon.Account {
	return &horizon.Account{
		AccountID: publicKey,
	}
}

func generateDRSAndFriends() {

	accountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Panic(err)
	}

	createDRSOp, drsKP := _operations.CreateAccountOp(&sourceAccount, "15")

	drsAccount := loadAccount(drsKP.Address())
	createTPListOp, tpListKP := _operations.CreateAccountOp(drsAccount, "1.5")

	createPFListOp, pfListKP := _operations.CreateAccountOp(drsAccount, "1.5")

	createREGListOp, regulatorListKP := _operations.CreateAccountOp(drsAccount, "2")

	createPriceSDGOp, priceSGDKP := _operations.CreateAccountOp(drsAccount, "1.5")

	createPriceTHBOp, priceTHBKP := _operations.CreateAccountOp(drsAccount, "1.5")

	createPriceUSDOp, priceUSDKP := _operations.CreateAccountOp(drsAccount, "1.5")

	createFirstRegulatorOp, regulatorKP := _operations.CreateAccountOp(drsAccount, "1")

	dropMasterWeight := txnbuild.Threshold(0)
	drsSignerWeight := txnbuild.Threshold(1)

	tpListAccount := loadAccount(tpListKP.Address())
	addDRSSignToTPListOp := _operations.SetSignerOp(tpListAccount, drsKP.Address(), drsSignerWeight, &dropMasterWeight, nil, nil, nil)

	regulatorListAccount := loadAccount(regulatorListKP.Address())
	addDRSSignToREGListOp := _operations.SetSignerOp(regulatorListAccount, drsKP.Address(), drsSignerWeight, &dropMasterWeight, nil, nil, nil)

	pfListAccount := loadAccount(pfListKP.Address())
	addDRSSignToPFListOp := _operations.SetSignerOp(pfListAccount, drsKP.Address(), drsSignerWeight, &dropMasterWeight, nil, nil, nil)

	lowThreshold := txnbuild.Threshold(254)
	mediumThreshold := txnbuild.Threshold(255)
	highThreshold := txnbuild.Threshold(254)

	drsSignerWeight = 254

	priceSGDAccount := loadAccount(priceSGDKP.Address())
	addDRSSignToPriceSGDOp := _operations.SetSignerOp(priceSGDAccount, drsKP.Address(), drsSignerWeight, &dropMasterWeight, &lowThreshold, &mediumThreshold, &highThreshold)

	priceTHBAccount := loadAccount(priceTHBKP.Address())
	addDRSSignToPriceTHBOp := _operations.SetSignerOp(priceTHBAccount, drsKP.Address(), drsSignerWeight, &dropMasterWeight, &lowThreshold, &mediumThreshold, &highThreshold)

	priceUSDAccount := loadAccount(priceUSDKP.Address())
	addDRSSignToPriceUSDOp := _operations.SetSignerOp(priceUSDAccount, drsKP.Address(), drsSignerWeight, &dropMasterWeight, &lowThreshold, &mediumThreshold, &highThreshold)

	drsManageDataTPListOp := _operations.ManageDataOp(drsAccount, "TrustedPartnerList", tpListKP.Address())

	drsManageDataREGListOp := _operations.ManageDataOp(drsAccount, "RegulatorList", regulatorListKP.Address())

	drsManageDataPFListOp := _operations.ManageDataOp(drsAccount, "PriceFeederList", pfListKP.Address())

	drsManageDataPriceUSDVELOOp := _operations.ManageDataOp(drsAccount, "Price[USD-VELO]", priceUSDKP.Address())

	drsManageDataPriceTHBVELOOp := _operations.ManageDataOp(drsAccount, "Price[THB-VELO]", priceTHBKP.Address())

	drsManageDataPriceSGDVELOOp := _operations.ManageDataOp(drsAccount, "Price[SGD-VELO]", priceSGDKP.Address())

	drsTrustLineVELOOp := _operations.ChangeTrustOp(drsAccount, veloAsset, veloAssetIssuerAccount)

	regulatorListManageDataAddRegulatorOp := _operations.ManageDataOp(regulatorListAccount, regulatorKP.Address(), "true")

	txFuture := txnbuild.Transaction{
		SourceAccount: &sourceAccount,
		Operations: []txnbuild.Operation{
			createDRSOp,
			createTPListOp,
			createPFListOp,
			createREGListOp,
			createPriceSDGOp,
			createPriceTHBOp,
			createPriceUSDOp,
			createFirstRegulatorOp,

			addDRSSignToTPListOp,
			addDRSSignToREGListOp,
			addDRSSignToPFListOp,
			addDRSSignToPriceSGDOp,
			addDRSSignToPriceTHBOp,
			addDRSSignToPriceUSDOp,

			drsManageDataTPListOp,
			drsManageDataREGListOp,
			drsManageDataPFListOp,
			drsManageDataPriceUSDVELOOp,
			drsManageDataPriceTHBVELOOp,
			drsManageDataPriceSGDVELOOp,
			drsTrustLineVELOOp,

			regulatorListManageDataAddRegulatorOp,
		},
		Timebounds: txnbuild.NewTimeout(300),
		Network:    network.TestNetworkPassphrase,
		BaseFee:    100 * 100,
	}

	log.Println("drs public:", drsKP.Address())
	log.Println("drs seed:", drsKP.Seed())

	log.Println("tpList public:", tpListKP.Address())
	log.Println("tpList seed:", tpListKP.Seed())

	log.Println("regList public:", regulatorListKP.Address())
	log.Println("regList seed:", regulatorListKP.Seed())

	log.Println("pfList public:", pfListKP.Address())
	log.Println("pfList seed:", pfListKP.Seed())

	log.Println("priceUSD public:", priceUSDKP.Address())
	log.Println("priceUSD seed:", priceUSDKP.Seed())

	log.Println("priceSGD public:", priceSGDKP.Address())
	log.Println("priceSGD seed:", priceSGDKP.Seed())

	log.Println("priceTHB public:", priceTHBKP.Address())
	log.Println("priceTHB seed:", priceTHBKP.Seed())

	log.Println("regulator public:", regulatorKP.Address())
	log.Println("regulator seed:", regulatorKP.Seed())

	txe, err := txFuture.BuildSignEncode(sourceKP, drsKP, tpListKP, pfListKP, regulatorListKP, priceSGDKP, priceTHBKP, priceUSDKP)
	if err != nil {
		log.Panic(err)
	}

	log.Println(txe)

}
