package main

import (
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	"log"
)

func main() {

	kp, err := vconvert.SecretKeyToKeyPair("SAAJAVEWZG5PGMYZVX377JZAP2N6DEDCKNEIL6HDHRNTOEZPGF4RM6EK")

	xdrTransaction := "AAAAACbhTnXxKBsclzxGrhp+2ZGU2WzzZ8w91VUjKkx029vqAAAAZAAN88cAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAJuFOdfEoGxyXPEauGn7ZkZTZbPNnzD3VVSMqTHTb2+oAAAAAAAAAAAX14QAAAAAAAAAAAA=="
	transaction, err := txnbuild.TransactionFromXDR(xdrTransaction)
	if err != nil {
		panic(err)
	}

	transaction.Network = network.TestNetworkPassphrase
	transaction.Timebounds = txnbuild.NewTimeout(300)
	txe, err := transaction.BuildSignEncode(kp)
	if err != nil {
		panic(err)
	}
	log.Println("singed transaction: ", txe)
}
