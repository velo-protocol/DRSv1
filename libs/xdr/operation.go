// Package vxdr contains the code for parsing the xdr structures to/from xdr string.
package vxdr

import (
	"fmt"
)

// VeloOp struct contains an operation body which contain an info on how the operation
// should be executed.
type VeloOp struct {
	Body OperationBody
}

// NewOperationBody creates a new OperationBody from a defined OperationType and its content.
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
	case OperationTypeRedeemCredit:
		tv, ok := value.(RedeemCreditOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be RedeemCreditOp")
		}
		opBody.RedeemCreditOp = &tv
	case OperationTypeRebalanceReserve:
		tv, ok := value.(RebalanceReserveOp)
		if !ok {
			return OperationBody{}, fmt.Errorf("invalid value, must be RebalanceReserveOp")
		}
		opBody.RebalanceReserveOp = &tv
	default:
		return OperationBody{}, fmt.Errorf("unknown operation type")
	}
	return opBody, nil
}
