package vxdr

import (
	"fmt"
)

type VeloOp struct {
	Body OperationBody
}

func NewOperationBody(opType OperationType, value interface{}) (OperationBody, error) {
	var opBody OperationBody
	opBody.Type = opType

	switch OperationType(opType) {
	case OperationTypeWhiteList:
		tv, ok := value.(WhiteListOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be WhiteListOp")
		}
		opBody.WhiteListOp = &tv
	default:
		return OperationBody{}, fmt.Errorf("unknown operation type")
	}
	// TODO: case OperationTypeSetupAccount
	return opBody, nil
}
