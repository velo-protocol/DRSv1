package main

import (
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
	veloTxB64, err := (&vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &vtxnbuild.WhiteList{
			Address: publicKey2,
			Role:    string(vxdr.RoleTrustedPartner),
		},
	}).BuildSignEncode(kp1)

	log.Printf("Error: %s", err)
	log.Printf("Velo tx xdr string: %s", veloTxB64)

	veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)

	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(veloTx.SourceAccount.GetAccountID())
	//log.Printf("Error: %s, %s", string(txSenderKeyPair.Hint()), veloTx.TxEnvelope().Signatures[0].Hint)
	log.Print(txSenderKeyPair.Hint() == veloTx.TxEnvelope().Signatures[0].Hint)
	if txSenderKeyPair.Hint() != veloTx.TxEnvelope().Signatures[0].Hint {
		log.Printf("Error: %s", err)
	}

}
