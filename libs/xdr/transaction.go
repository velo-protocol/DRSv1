package vxdr

import (
	"github.com/stellar/go/xdr"
)

// VeloTx is a struct that contains a VeloOp and a SourceAccount.
type VeloTx struct {
	SourceAccount xdr.AccountId
	VeloOp        VeloOp
}

// VeloTxEnvelope is struct that wraps VeloTx and its signatures.
type VeloTxEnvelope struct {
	VeloTx     VeloTx
	Signatures []xdr.DecoratedSignature `xdrmaxsize:"20"`
}
