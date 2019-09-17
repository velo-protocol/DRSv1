package vxdr

import (
	"github.com/stellar/go/xdr"
)

type OperationType int32

const (
	OperationTypeWhiteList OperationType = 0
	// TODO: OperationTypeSetupAccount OperationType = 1
)

type OperationBody struct {
	Type        OperationType
	WhiteListOp *WhiteListOp
	// TODO: SetupAccountOp *SetupAccountOp
}

type WhiteListOp struct {
	Address xdr.AccountId
	Role    Role
}
