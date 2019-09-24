package helper

import (
	"encoding/json"
	"fmt"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
)

const (
	PublicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
	SecretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
	PublicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
	SecretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
)

var (
	KP1, _ = vconvert.SecretKeyToKeyPair(SecretKey1)
	KP2, _ = vconvert.SecretKeyToKeyPair(SecretKey2)
)

func DecodeB64VeloTx(base64VeloTx string) {
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

func CompareVeloTxSigner(base64VeloTx string, accountToCompare string) {
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
