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
	buildB64WhitelistOp(helper.PublicKey1, helper.PublicKey2, vxdr.RoleTrustedPartner, helper.KP1)

	helper.DecodeB64VeloTx("AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAQAAAAAAAAABAAAAAACYloAAAAADVEhCAAAAAAR2VEhCAAAAAdbT7LsAAABAVmoek8shsDnBLATJupu2ACmVrk8olEj+r3QOhY78ARvIXZA1F9Se5hw1/GArw8q9sI3JxR521ZEDQzBIZj0HAg==")
	helper.CompareVeloTxSigner("AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAQAAAAAAAAABAAAAAACYloAAAAADVEhCAAAAAAR2VEhCAAAAAdbT7LsAAABAVmoek8shsDnBLATJupu2ACmVrk8olEj+r3QOhY78ARvIXZA1F9Se5hw1/GArw8q9sI3JxR521ZEDQzBIZj0HAg==", helper.PublicKey1)

}

func buildB64WhitelistOp(txSourceAccount, opSourceAccount string, whiteListRole vxdr.Role, secretKey *keypair.Full) {
	fmt.Println("##### Start Build WhiteList Operation #####")

	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: txSourceAccount,
		},
		VeloOp: &vtxnbuild.WhiteList{
			Address: opSourceAccount,
			Role:    string(whiteListRole),
		},
	}).BuildSignEncode(secretKey)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Velo Transaction: %s \n", veloTxB64)

	fmt.Println("##### End Build WhiteList Operation #####")
}
