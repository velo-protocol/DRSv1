package main

import (
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/convert"
	_operations "github.com/velo-protocol/DRSv1/tools/shared/operations"
	"log"
)

var (
	sourceSeedKey          = "<YOUR_SECRET_KEY>"
	veloAssetIssuerSeedKey = "<VELO_ASSET_ISSUER_SEED_KEY>"
	veloAsset              = "VELO"

	sourceKP, _          = vconvert.SecretKeyToKeyPair(sourceSeedKey)
	veloAssetIssuerKP, _ = vconvert.SecretKeyToKeyPair(veloAssetIssuerSeedKey)
	client               = horizonclient.DefaultTestNetClient
)

func main() {
	mintVeloToDistributor()
}

func loadAccount(publicKey string) *horizon.Account {
	return &horizon.Account{
		AccountID: publicKey,
	}
}

func mintVeloToDistributor() {

	accountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Panic(err)
	}

	veloIssuerAccount := loadAccount(veloAssetIssuerKP.Address())

	// Create VELO Distributor Account
	createVeloDistributorOp, veloDistributorKP := _operations.CreateAccountOp(&sourceAccount, "1.5")

	// Add Operation ChangeTrust VELO
	veloDistributorAccount := loadAccount(veloDistributorKP.Address())
	veloAssetIssuerTrustLineVELOOp := _operations.ChangeTrustOp(veloDistributorAccount, veloAsset, veloAssetIssuerKP.Address())

	// Add Operation Payment To ISSUED 900,000,000.0000000 VELO Onchain
	veloIssuedAmount := "900000000"
	paymentOp := _operations.PaymentOp(veloIssuerAccount, veloAsset, veloAssetIssuerKP.Address(), veloDistributorKP.Address(), veloIssuedAmount)

	// Add Operation SetSigner Drop VELO Issuer Master Key
	veloIssuerSignerWeight := txnbuild.Threshold(0)
	dropVeloIssuerOp := &txnbuild.SetOptions{
		MasterWeight:  &veloIssuerSignerWeight,
		SourceAccount: veloIssuerAccount,
	}

	tx := txnbuild.Transaction{
		SourceAccount: &sourceAccount,
		Operations: []txnbuild.Operation{
			createVeloDistributorOp,

			veloAssetIssuerTrustLineVELOOp,

			paymentOp,

			dropVeloIssuerOp,
		},
		Timebounds: txnbuild.NewTimeout(300),
		Network:    network.TestNetworkPassphrase,
	}

	log.Println("VELO Issuer Address: ", veloAssetIssuerKP.Address())
	log.Println("VELO Issuer Seed: ", veloAssetIssuerKP.Seed())

	log.Println("VELO Distributor Address: ", veloDistributorKP.Address())
	log.Println("VELO Distributor Seed: ", veloDistributorKP.Seed())

	txe, err := tx.BuildSignEncode(sourceKP, veloDistributorKP, veloAssetIssuerKP)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Transaction XDR: ", txe)

	txSuccess, err := client.SubmitTransactionXDR(txe)
	if err != nil {
		herr, ok := err.(horizonclient.Error)
		if !ok {
			log.Panic("fail to confirm with stellar")
		}
		herrString, _ := herr.ResultString()
		log.Panic(err, fmt.Sprintf(`horizon err "%s"`, herrString))
	}

	log.Println("Transaction Hash: ", txSuccess)
}
