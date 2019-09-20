package main

import (
	"encoding/json"
	"fmt"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"log"
)

const (
	publicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
	secretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
	publicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
	secretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
)

var (
	kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)
)

func main() {
	// Exists roles
	// vxdr.RoleTrustedPartner Role = "TRUSTED_PARTNER"
	// vxdr.RolePriceFeeder    Role = "PRICE_FEEDER"
	// vxdr.RoleRegulator      Role = "REGULATOR"
	buildB64WhitelistOp(publicKey1, publicKey2, vxdr.RoleTrustedPartner, kp1)

	decodeB64VeloTx("AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAAAAAAEAAAAAtRdjMD9LOpHwyX4gB3qQki+O1VIWOgJ9CKnXfORHzloAAAAJUkVHVUxBVE9SAAAAAAAAAdbT7LsAAABAK1dEk1kZUbZ2ORyAsfLqmoE6XGaBaB41vy18udY95bnhg58+n5FrRrOJwzmWmW86qLhSJ0ZwocLjQ2JevmbEAA==")

	compareVeloTxSigner("AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAAAAAAEAAAAAtRdjMD9LOpHwyX4gB3qQki+O1VIWOgJ9CKnXfORHzloAAAAJUkVHVUxBVE9SAAAAAAAAAdbT7LsAAABAK1dEk1kZUbZ2ORyAsfLqmoE6XGaBaB41vy18udY95bnhg58+n5FrRrOJwzmWmW86qLhSJ0ZwocLjQ2JevmbEAA==", publicKey1)
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

func buildB64SetupCreditOp() {
	// TODO: Add more
}

func decodeB64VeloTx(base64VeloTx string) {
	fmt.Println("##### Start Decode Base64 Velo Transaction #####")
	veloTx, err := vtxnbuild.TransactionFromXDR(base64VeloTx)
	if err != nil {
		panic(err)
	}
	veloTxByte, err := json.Marshal(veloTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Velo Transaction: %s \n", string(veloTxByte))

	fmt.Println("##### End Decode Base64 Velo Transaction #####")
}

func compareVeloTxSigner(base64VeloTx string, accountToCompare string) {
	fmt.Println("##### Start Compare Velo Transaction Signer #####")
	veloTx, err := vtxnbuild.TransactionFromXDR(base64VeloTx)
	if err != nil {
		panic(err)
	}
	compareAccount, err := vconvert.PublicKeyToKeyPair(accountToCompare)
	if err != nil {
		panic(err)
	}

	isMatch := veloTx.TxEnvelope().Signatures[0].Hint == compareAccount.Hint()

	fmt.Printf("Velo Tx Signer Account == %s : %+v \n", compareAccount.Address(), isMatch)

	fmt.Println("##### End Compare Velo Transaction Signer #####")
}
