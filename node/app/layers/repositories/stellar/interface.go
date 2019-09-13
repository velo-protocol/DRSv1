package stellar

import "github.com/stellar/go/protocols/horizon"

type Repo interface {
	BuildSetupTx(drsAccount *horizon.Account, peggedValue string, peggedCurrency string, assetName string, creditOwnerAddress string) (setupTxB64 string, issuerAddress string, distributorAddress string, err error)
	BuildMintTx(drsAccount *horizon.Account, amount string, assetName string, issuerAddress string, distributorAddress string) (string, error)
	LoadAccount(stellarAddress string) (*horizon.Account, error)
	SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error)
}
