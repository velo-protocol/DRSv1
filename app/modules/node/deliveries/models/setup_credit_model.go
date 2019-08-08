package models

import "github.com/mitchellh/mapstructure"

type SetupCreditRequest struct {
	SignedIssuerCreationTx string `json:"signedIssuerCreationTx" binding:"required"`
	PeggedValue            string `json:"peggedValue" binding:"required"`
	PeggedCurrency         string `json:"peggedCurrency" binding:"required"`
	AssetName              string `json:"assetName" binding:"required"`
}

type SetupCreditResponse struct {
	CreditOwnerAddress   string `json:"creditOwnerAddress"`
	IssuerCreationTxHash string `json:"issuerCreationTxHash"`
	SetupTxHash          string `json:"setupTxHash"`
	IssuerAddress        string `json:"issuerAddress"`
	DistributorAddress   string `json:"distributorAddress"`
	AssetName            string `json:"assetName"`
	PeggedValue          string `json:"peggedValue"`
	PeggedCurrency       string `json:"peggedCurrency"`
}

func (model *SetupCreditResponse) Parse(data interface{}) (*SetupCreditResponse, error) {
	err := mapstructure.Decode(data, &model)
	return model, err
}
