package operation

type Interface interface {
	Setup(
		peggedValue string,
		peggedCurrency string,
		assetName string,
		creditOwnerAddress string,
	) (setupTxB64 string, issuerAddress string, distributorAddress string, err error)

	Mint(
		amount string,
		assetName string,
		issuerAddress string,
		distributorAddress string,
	) (string, error)
}
