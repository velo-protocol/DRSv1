package main

import (
	"fmt"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/testkit/helper"
	"log"
)

func main() {
	veloTxB64 := buildB64WhitelistOp(helper.PublicKeyFirstRegulator, helper.PublicKeyPF, vxdr.RolePriceFeeder, "THB", helper.KPFirstRegulator)

	helper.DecodeB64VeloTx(veloTxB64)
	helper.CompareVeloTxSigner(veloTxB64, helper.PublicKeyFirstRegulator)
}

func buildB64WhitelistOp(txSourceAccount, opSourceAccount string, whiteListRole vxdr.Role, currency string, secretKey *keypair.Full) string {
	fmt.Println("##### Start Build WhiteList Operation #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.WhiteList{
			Address:  opSourceAccount,
			Role:     string(whiteListRole),
			Currency: currency,
		},
	}).BuildSignEncode(secretKey)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build WhiteList Operation #####")

	return veloTxB64
}
