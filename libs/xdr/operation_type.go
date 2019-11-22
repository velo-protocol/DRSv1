package vxdr

import (
	"github.com/stellar/go/xdr"
)

// OperationType is an enum that indicate the operation type.
type OperationType int32

const (
	OperationTypeWhitelist        OperationType = 0
	OperationTypeSetupCredit      OperationType = 1
	OperationTypePriceUpdate      OperationType = 2
	OperationTypeMintCredit       OperationType = 3
	OperationTypeRedeemCredit     OperationType = 4
	OperationTypeRebalanceReserve OperationType = 5
)

// OperationBody is an struct that contains the data of each op type.
type OperationBody struct {
	Type               OperationType
	WhitelistOp        *WhitelistOp
	SetupCreditOp      *SetupCreditOp
	PriceUpdateOp      *PriceUpdateOp
	MintCreditOp       *MintCreditOp
	RedeemCreditOp     *RedeemCreditOp
	RebalanceReserveOp *RebalanceReserveOp
}

// WhitelistOp is an struct that contains parameter required to perform Whitelist.
type WhitelistOp struct {
	Address  xdr.AccountId
	Role     Role
	Currency Currency
}

// SetupCreditOp is an struct that contains parameter required to perform SetupCredit.
type SetupCreditOp struct {
	PeggedValue    xdr.Int64
	PeggedCurrency string
	AssetCode      string
}

// SetupCreditOp is an struct that contains parameter required to perform SetupCredit.
type PriceUpdateOp struct {
	Asset                       string
	Currency                    Currency
	PriceInCurrencyPerAssetUnit xdr.Int64
}

// MintCreditOp is an struct that contains parameter required to perform MintCredit.
type MintCreditOp struct {
	AssetCodeToBeIssued string
	CollateralAssetCode Asset
	CollateralAmount    xdr.Int64
}

// RedeemCreditOp is an struct that contains parameter required to perform RedeemCredit.
type RedeemCreditOp struct {
	AssetCode string
	Issuer    xdr.AccountId
	Amount    xdr.Int64
}

// RebalanceReserveOp is an struct that contains parameter required to perform RebalanceReserve.
type RebalanceReserveOp struct {
}
