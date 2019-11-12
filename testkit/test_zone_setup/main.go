package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/testkit/helper"
	"log"
	"net/http"
)

var tp, pfTHB, pfUSD, pfSGD, redeemer, source *keypair.Full
var tpClient, regClient *vclient.Client

func main() {
	var err error

	// Create and fund account
	tp = createAccount()
	log.Printf("TP: %s / %s", tp.Address(), tp.Seed())

	pfTHB = createAccount()
	log.Printf("PF THB: %s / %s", pfTHB.Address(), pfTHB.Seed())

	pfUSD = createAccount()
	log.Printf("PF USD: %s / %s", pfUSD.Address(), pfUSD.Seed())

	pfSGD = createAccount()
	log.Printf("PF SGD: %s / %s", pfSGD.Address(), pfSGD.Seed())

	redeemer = createAccount()
	log.Printf("Redeemer: %s / %s", redeemer.Address(), redeemer.Seed())

	source = createAccount() // for tx submission

	// Whitelisting TP and PFs by a regulator
	regClient, err = vclient.NewDefaultTestNetClient(helper.VeloNodeUrl, helper.SecretKeyFirstRegulator)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = regClient.Close()
	}()

	whitelistTP(tp)
	whitelistPF(pfTHB, "THB")
	whitelistPF(pfUSD, "USD")
	whitelistPF(pfSGD, "SGD")

	// Setting up token by TP
	tpClient, err = vclient.NewDefaultTestNetClient(helper.VeloNodeUrl, tp.Seed())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = tpClient.Close()
	}()

	vTHB := setupCredit("THB")
	vUSD := setupCredit("USD")
	vSGD := setupCredit("SGD")

	// Source account
	sourceAccount, err := horizonclient.DefaultTestNetClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: source.Address(),
	})
	if err != nil {
		panic(err)
	}

	tx := txnbuild.Transaction{
		SourceAccount: &sourceAccount,
		Operations: []txnbuild.Operation{
			// TODO: add trust line of virtual credit to redeemer
			&txnbuild.ChangeTrust{
				SourceAccount: &horizon.Account{
					AccountID: redeemer.Address(),
				},
				Line: vTHB,
			},
			&txnbuild.ChangeTrust{
				SourceAccount: &horizon.Account{
					AccountID: redeemer.Address(),
				},
				Line: vUSD,
			},
			&txnbuild.ChangeTrust{
				SourceAccount: &horizon.Account{
					AccountID: redeemer.Address(),
				},
				Line: vSGD,
			},
			// TODO: add trust line of VELO to redeemer
			&txnbuild.ChangeTrust{
				SourceAccount: &horizon.Account{
					AccountID: redeemer.Address(),
				},
				Line: txnbuild.CreditAsset{
					Code:   helper.VeloAssetCode,
					Issuer: helper.VeloIssuerAddress,
				},
			},
			// TODO: add trust line of VELO to TP
			&txnbuild.ChangeTrust{
				SourceAccount: &horizon.Account{
					AccountID: tp.Address(),
				},
				Line: txnbuild.CreditAsset{
					Code:   helper.VeloAssetCode,
					Issuer: helper.VeloIssuerAddress,
				},
			},
		},
		Timebounds: txnbuild.NewInfiniteTimeout(),
		Network:    network.TestNetworkPassphrase,
	}

	signedTx, err := tx.BuildSignEncode(source, tp, redeemer)
	if err != nil {
		panic(err)
	}

	_, err = horizonclient.DefaultTestNetClient.SubmitTransactionXDR(signedTx)
	if err != nil {
		panic(err)
	}

	log.Println("Done!")
}

func createAccount() *keypair.Full {
	kp, err := keypair.Random()
	if err != nil {
		panic(errors.Wrap(err, "failed to generate keypair"))
	}

	resp, err := http.Get(fmt.Sprintf("https://horizon-testnet.stellar.org/friendbot?addr=%s", kp.Address()))
	if err != nil {
		panic(errors.Wrap(err, "failed to get free lumens from friendbot"))
	}

	if resp.StatusCode != http.StatusOK {
		panic(errors.New("failed to get free lumens from friendbot"))
	}

	return kp
}

func whitelistTP(kp *keypair.Full) {
	_, err := regClient.Whitelist(context.Background(), vtxnbuild.Whitelist{
		Address: kp.Address(),
		Role:    string(vxdr.RoleTrustedPartner),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("TP has been whitelisted: %s / %s", kp.Address(), vxdr.RoleTrustedPartner)
}

func whitelistPF(kp *keypair.Full, currency string) {
	_, err := regClient.Whitelist(context.Background(), vtxnbuild.Whitelist{
		Address:  kp.Address(),
		Role:     string(vxdr.RolePriceFeeder),
		Currency: currency,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("PF has been whitelisted: %s / %s / %s", kp.Address(), vxdr.RolePriceFeeder, currency)
}

func setupCredit(currency string) *txnbuild.CreditAsset {
	result, err := tpClient.SetupCredit(context.Background(), vtxnbuild.SetupCredit{
		PeggedValue:    "1",
		PeggedCurrency: currency,
		AssetCode:      "v" + currency,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Credit has been setup: %s / %s", result.VeloNodeResult.AssetCode, result.VeloNodeResult.AssetIssuer)

	return &txnbuild.CreditAsset{
		Code:   result.VeloNodeResult.AssetCode,
		Issuer: result.VeloNodeResult.AssetIssuer,
	}
}
