package services

type Operation interface {
	Setup(
		peggedValue string,
		peggedCurrency string,
		assetName string,
		creditOwnerAddress string,
	) (string, error)

	Mint(
		amount string,
		assetName string,
		issuerAddress string,
		distributorAddress string,
	) (string, error)
}
