package vxdr

import (
	"github.com/stellar/go/xdr"
)

type VeloTx struct {
	SourceAccount xdr.AccountId
	VeloOp        VeloOp
}

type VeloTxEnvelope struct {
	VeloTx     VeloTx
	Signatures []xdr.DecoratedSignature `xdrmaxsize:"20"`
}