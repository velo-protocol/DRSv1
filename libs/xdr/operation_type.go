package vxdr

import (
	"github.com/stellar/go/xdr"
)

type OperationType int32

const (
	OperationTypeWhiteList   OperationType = 0
	OperationTypeSetupCredit OperationType = 1
)

type OperationBody struct {
	Type          OperationType
	WhiteListOp   *WhiteListOp
	SetupCreditOp *SetupCreditOp
}

type WhiteListOp struct {
	Address xdr.AccountId
	Role    Role
}

type SetupCreditOp struct {
	PeggedValue    xdr.Int64
	PeggedCurrency string
	AssetName      string
}
