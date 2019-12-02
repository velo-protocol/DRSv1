package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"log"
)

var (
	sourceSeedKey   = "<YOUR_SECRET_KEY>"
	veloAssetIssuer = "<VELO_ISSUER_PUBLIC_KEY>"
	veloAsset       = "VELO"
	sourceKP, _     = vconvert.SecretKeyToKeyPair(sourceSeedKey)
	client          = horizonclient.DefaultTestNetClient
)

func main() {
	accountRequest := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Panic(err)
	}
	sugarDaddyFaucetAccountKP, err := keypair.Random()
	if err != nil {
		panic(err)
	}

	log.Printf("Sugar daddy Faucet: %s / %s", sugarDaddyFaucetAccountKP.Address(), sugarDaddyFaucetAccountKP.Seed())

	tx := txnbuild.Transaction{
		SourceAccount: &sourceAccount,
		Operations: []txnbuild.Operation{
			&txnbuild.CreateAccount{
				Destination:   sugarDaddyFaucetAccountKP.Address(),
				Amount:        "1.5",
				SourceAccount: &sourceAccount,
			},
			&txnbuild.ChangeTrust{
				Line: txnbuild.CreditAsset{
					Code:   veloAsset,
					Issuer: veloAssetIssuer,
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: sugarDaddyFaucetAccountKP.Address(),
				},
			},
		},
		Timebounds: txnbuild.NewTimeout(300),
		Network:    network.TestNetworkPassphrase,
	}

	signedTx, err := tx.BuildSignEncode(sourceKP, sugarDaddyFaucetAccountKP)
	if err != nil {
		log.Panic(err)
	}

	_, err = client.SubmitTransactionXDR(signedTx)
	if err != nil {
		panic(err)
	}

	log.Println("Done!")

}
