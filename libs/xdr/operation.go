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

	switch opType {
	case OperationTypeWhitelist:
		tv, ok := value.(WhitelistOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be WhitelistOp")
		}
		opBody.WhitelistOp = &tv
	case OperationTypeSetupCredit:
		tv, ok := value.(SetupCreditOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be SetupCreditOp")
		}
		opBody.SetupCreditOp = &tv
	case OperationTypePriceUpdate:
		tv, ok := value.(PriceUpdateOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be PriceUpdateOp")
		}
		opBody.PriceUpdateOp = &tv
	case OperationTypeMintCredit:
		tv, ok := value.(MintCreditOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be MintCreditOp")
		}
		opBody.MintCreditOp = &tv
	default:
		return OperationBody{}, fmt.Errorf("unknown operation type")
	}
	return opBody, nil
}
