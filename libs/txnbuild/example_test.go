package vtxnbuild

import (
	"github.com/stellar/go/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"log"
)

func ExampleVeloTx_Build() {
	veloTx := VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &Whitelist{
			Address: publicKey2,
			Role:    string(vxdr.RoleTrustedPartner),
		},
	}

	err := veloTx.Build()
	if err != nil {
		log.Println(err)
		return
	}
}

func ExampleVeloTx_Sign() {

	veloTx := VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &Whitelist{
			Address: publicKey2,
			Role:    string(vxdr.RoleTrustedPartner),
		},
	}
	_ = veloTx.Build()
	err := veloTx.Sign(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
}

func ExampleVeloTx_Base64() {
	veloTx := VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &Whitelist{
			Address: publicKey2,
			Role:    string(vxdr.RoleTrustedPartner),
		},
	}
	_ = veloTx.Build()
	veloTxB64, err := veloTx.Base64()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAAAAAAEAAAAAtRdjMD9LOpHwyX4gB3qQki+O1VIWOgJ9CKnXfORHzloAAAAPVFJVU1RFRF9QQVJUTkVSAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
}

func ExampleVeloTx_BuildSignEncode() {
	veloTxB64, err := (&VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: publicKey1,
		},
		VeloOp: &Whitelist{
			Address: publicKey2,
			Role:    string(vxdr.RoleTrustedPartner),
		},
	}).BuildSignEncode(kp1, kp2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(veloTxB64)
	// AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAAAAAAAAEAAAAAtRdjMD9LOpHwyX4gB3qQki+O1VIWOgJ9CKnXfORHzloAAAAPVFJVU1RFRF9QQVJUTkVSAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALW0+y7AAAAQEQI7/oYaybfF0m5jQDO72fFpzx1YyyiJn7keTGp34B6hpgbCobN8M41WERPEP+Z+kRSxDkzJe49GT7jicQhTwLkR85aAAAAQFKPYrtXStZiH4mDSIcmke1UhwkP6URJ8JODgBcRRGpMaPFxiH/mEA6sxuu/+TvFf6ZRcnb6twBd3yRU2hDhNQk=
}
