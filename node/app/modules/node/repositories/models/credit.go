package models

import (
	"github.com/mitchellh/mapstructure"
)

type Credit struct {
	CreditOwnerAddress   string `json:"creditOwnerAddress"`
	IssuerCreationTxHash string `json:"issuerCreationTxHash"`
	SetupTxHash          string `json:"setupTxHash"`
	IssuerAddress        string `json:"issuerAddress"`
	DistributorAddress   string `json:"distributorAddress"`
	AssetName            string `json:"assetName"`
	PeggedValue          string `json:"peggedValue"`
	PeggedCurrency       string `json:"peggedCurrency"`
}

func (model *Credit) Parse(data interface{}) (*Credit, error) {
	err := mapstructure.Decode(data, &model)
	return model, err
}
