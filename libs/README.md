# Velo Libs

## Velo Transaction Builder
Use `cen/libs/txnbuild` when you want to build a Velo Transaction and convert to an XDR format. The package usage is designed to be similar to `stellar/go/txnbuild`.

### Usage
```go
package main

import (
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/vxdr"
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
	// Create Velo Transaction
    veloTx := vtxnbuild.VeloTx{
        SourceAccount: &txnbuild.SimpleAccount{
            AccountID: "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73",
        },
        VeloOp: &vtxnbuild.WhiteList{
            Address: "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25",
            Role:    string(vxdr.RoleTrustedPartner),
        },
    }
    
    // Build and sign the transaction
    veloXdrTxB64, err := veloTx.BuildSignEncode(kp1, kp2)
    if err != nil {
    	panic(err)
    }
    
    // ...
    // submit to velo node
    // ...
    
    // Decode XDR
    newVeloTx, err := vtxnbuild.TransactionFromXDR(veloXdrTxB64)
    
}
```
