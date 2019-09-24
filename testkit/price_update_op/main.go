package main

import (
	"fmt"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/testkit/helper"
	"log"
)

func main() {
	veloTxB64 := buildB64PriceUpdateOp(helper.PublicKeyTP, "VELO", "THB", "1", helper.KPTP)

	helper.DecodeB64VeloTx(veloTxB64)
	helper.CompareVeloTxSigner(veloTxB64, helper.PublicKeyRegulator)
}

func buildB64PriceUpdateOp(txSourceAccount, asset, currency, priceInCurrencyPerAssetUnit string, secretKey *keypair.Full) string {
	fmt.Println("##### Start Build Price Update Operation #####")

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

	fmt.Println("##### End Build Price Update Operation #####")

	return veloTxB64
}
