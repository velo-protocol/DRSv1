package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"log"
)

var (
	sourceSeedKey    = "<YOUR_SECRET_KEY>"
	veloDistributor  = "<VELO_DISTRIBUTOR_SEED_KEY>"
	veloAssetIssuer  = "<VELO_ISSUER_PUBLIC_KEY>"
	veloAsset        = "VELO"
	drsReserve       = "<DRS_RESERVE_PUBLIC_KEY>"
	sugarDaddyFaucet = "<SUGAR_DADDY_FAUCET>_PUBLIC_KEY"

	sourceKP, _          = vconvert.SecretKeyToKeyPair(sourceSeedKey)
	veloDistributorKP, _ = vconvert.SecretKeyToKeyPair(veloDistributor)
	client               = horizonclient.DefaultTestNetClient
)

func main() {

	accountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Panic(err)
	}

	tx := txnbuild.Transaction{
		SourceAccount: &sourceAccount,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: drsReserve,
				Amount:      "500000000",
				Asset: txnbuild.CreditAsset{
					Code:   veloAsset,
					Issuer: veloAssetIssuer,
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: veloDistributorKP.Address(),
				},
			},
			&txnbuild.Payment{
				Destination: sugarDaddyFaucet,
				Amount:      "100000000",
				Asset: txnbuild.CreditAsset{
					Code:   veloAsset,
					Issuer: veloAssetIssuer,
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: veloDistributorKP.Address(),
				},
			},
		},
		Timebounds: txnbuild.NewTimeout(300),
		Network:    network.TestNetworkPassphrase,
	}

	signedTx, err := tx.BuildSignEncode(sourceKP, veloDistributorKP)
	if err != nil {
		log.Panic(err)
	}

	_, err = client.SubmitTransactionXDR(signedTx)
	if err != nil {
		panic(err)
	}

	log.Println("Done!")

}
