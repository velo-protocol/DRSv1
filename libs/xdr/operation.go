package vxdr

import (
	"fmt"
)

type VeloOp struct {
	Body OperationBody
}

func NewOperationBody(opType OperationType, value interface{}) (result OperationBody, err error) {
	result.Type = opType
	switch OperationType(opType) {
	case OperationTypeWhiteList:
		tv, ok := value.(WhiteListOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be WhiteListOp")
			return
		}
		result.WhiteListOp = &tv
	}
	// TODO: case OperationTypeSetupAccount
	return
}
