package vxdr

import (
	"github.com/stellar/go/xdr"
)

type OperationType int32

const (
	OperationTypeWhitelist    OperationType = 0
	OperationTypeSetupCredit  OperationType = 1
	OperationTypePriceUpdate  OperationType = 2
	OperationTypeMintCredit   OperationType = 3
	OperationTypeRedeemCredit OperationType = 4
)

type OperationBody struct {
	Type           OperationType
	WhitelistOp    *WhitelistOp
	SetupCreditOp  *SetupCreditOp
	PriceUpdateOp  *PriceUpdateOp
	MintCreditOp   *MintCreditOp
	RedeemCreditOp *RedeemCreditOp
}

type WhitelistOp struct {
	Address  xdr.AccountId
	Role     Role
	Currency Currency
}

type SetupCreditOp struct {
	PeggedValue    xdr.Int64
	PeggedCurrency string
	AssetCode      string
}

type PriceUpdateOp struct {
	Asset                       string
	Currency                    Currency
	PriceInCurrencyPerAssetUnit xdr.Int64
}

type MintCreditOp struct {
	AssetCodeToBeIssued string
	CollateralAssetCode Asset
	CollateralAmount    xdr.Int64
}

type RedeemCreditOp struct {
	AssetCode string
	Issuer    xdr.AccountId
	Amount    xdr.Int64
}
