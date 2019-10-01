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
	veloTxB64 := buildB64SetupCreditOp(helper.PublicKeyTP, "THB", "1", "vTHB", helper.KPTP)

	helper.DecodeB64VeloTx(veloTxB64)
	helper.CompareVeloTxSigner(veloTxB64, helper.PublicKeyFirstRegulator)

}

func buildB64SetupCreditOp(txSourceAccount, peggedCurrency, peggedValue, assetCode string, secretKey *keypair.Full) string {
	fmt.Println("##### Start Build Setup Credit Operation #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.SetupCredit{
			PeggedValue:    peggedValue,
			PeggedCurrency: peggedCurrency,
			AssetCode:      assetCode,
		},
	}).BuildSignEncode(secretKey)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build Setup Credit Operation #####")

	return veloTxB64
}
