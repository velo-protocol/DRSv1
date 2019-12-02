package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"log"
	"net/http"
)

var pfTHB, pfUSD, pfSGD *keypair.Full
var regClient *vclient.Client

var (
	veloNodeUrl      = "dev-velo-cen-node-01.velo.org:8080"
	regulatorSeedKey = "<SECRET_KEY>"

	regulatorKP, _ = vconvert.SecretKeyToKeyPair(regulatorSeedKey)
)

func main() {
	var err error

	// Create and fund account
	pfTHB = createAccount()
	log.Printf("PF THB: %s / %s", pfTHB.Address(), pfTHB.Seed())

	pfUSD = createAccount()
	log.Printf("PF USD: %s / %s", pfUSD.Address(), pfUSD.Seed())

	pfSGD = createAccount()
	log.Printf("PF SGD: %s / %s", pfSGD.Address(), pfSGD.Seed())

	// Whitelisting TP and PFs by a regulator
	regClient, err = vclient.NewDefaultTestNetClient(veloNodeUrl, regulatorKP.Seed())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = regClient.Close()
	}()

	whitelistPF(pfTHB, "THB")
	whitelistPF(pfUSD, "USD")
	whitelistPF(pfSGD, "SGD")
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
