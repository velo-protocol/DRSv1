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
	buildB64SetupCreditOp(helper.PublicKey1, "THB", "1", "vTHB", helper.KP1)

	helper.DecodeB64VeloTx("AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAQAAAAAAAAABAAAAAACYloAAAAADVEhCAAAAAAR2VEhCAAAAAdbT7LsAAABAVmoek8shsDnBLATJupu2ACmVrk8olEj+r3QOhY78ARvIXZA1F9Se5hw1/GArw8q9sI3JxR521ZEDQzBIZj0HAg==")
	helper.CompareVeloTxSigner("AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAQAAAAAAAAABAAAAAACYloAAAAADVEhCAAAAAAR2VEhCAAAAAdbT7LsAAABAVmoek8shsDnBLATJupu2ACmVrk8olEj+r3QOhY78ARvIXZA1F9Se5hw1/GArw8q9sI3JxR521ZEDQzBIZj0HAg==", helper.PublicKey1)

}

func buildB64SetupCreditOp(txSourceAccount, peggedCurrency, peggedValue, assetCode string, secretKey *keypair.Full) {
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
}
