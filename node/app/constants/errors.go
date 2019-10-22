package constants

var (
	ErrSignatureNotFound        = "signature not found"
	ErrUnknownVeloOperationType = "unknown velo operation type"
	ErrFormatMissingOperation   = "operation type %s is missing"
	ErrMustNotBeBlank           = "must not be blank"
	ErrInvalidFormat            = "invalid %s format"

	ErrSignatureNotMatchSourceAccount = "the signature and source account does not match"
	ErrFormatSignerNotHavePermission  = "the signer is not found or does not have sufficient permission to perform %s"
	ErrWhitelistAlreadyWhitelisted    = "the address %s has already been whitelisted for the role %s"

	ErrPriceFeederCurrencyMustNotBlank         = "currency must not be blank for price feeder role"
	ErrCurrencyMustBeBlank                     = "currency must be blank for non-price feeder role"
	ErrCurrencyMustMatchWithRegisteredCurrency = "currency must match with the registered currency"
	ErrGetSenderAccount                        = "fail to get tx sender account"
	ErrGetRoleListAccount                      = "fail to get role list accounts"

	ErrGetDrsAccountData                = "fail to get data of drs account"
	ErrGetTrustedPartnerListAccountData = "fail to get data of trusted partner list account"
	ErrGetPriceFeederListAccountData    = "fail to get data of price feeder list account"

	ErrGetAccountDetail                   = "fail to get account detail of %s"
	ErrGetDrsAccountDetail                = "fail to get account detail of drs account"
	ErrGetTrustedPartnerAccountDetail     = "fail to get account detail of trusted partner account"
	ErrGetTrustedPartnerMetaAccountDetail = "fail to get account detail of trusted partner's meta account"
	ErrGetIssuerAccountDetail             = "fail to get account detail of issuer account"
	ErrGetIssuerAccount                   = "fail to get issuer account"
	ErrGetAsset                           = "fail to get asset %s"

	ErrVerifyTrustedPartnerAccount = "fail to verify trusted partner account"
	ErrVerifyAssetCode             = "fail to verify asset code"

	ErrDerivedKeyPairFromSeed   = "fail to derived KP from seed key"
	ErrBuildAndSignTransaction  = "fail to build and sign tx"
	ErrUnknownRoleSpecified     = "unknown role specified"
	ErrAssetCodeAlreadyBeenUsed = "asset code %s has already been used"
	ErrToDecodeData             = `fail to decode data "%s"`
	ErrKeyNotFoundInAccountData = `key "%s" not found in the data of account %s`
	ErrNoValidPrice             = "no valid price in price account %s"

	ErrCreateTrustedPartnerMetaKeyPair = "fail to create trusted partner meta KP"
	ErrCreateIssuerKeyPair             = "fail to create issuer KP"
	ErrCreateDistributorKeyPair        = "fail to create distributor KP"

	ErrAssetCodeToBeIssuedNotSetup      = "assetCodeToBeIssued %s have not been setup yet"
	ErrGetPriceOfPeggedCurrency         = "fail to get price of pegged currency"
	ErrPeggedValueMustBeGreaterThanZero = "pegged value must be greater than zero"
	ErrMedianPriceMustBeGreaterThanZero = "median price must be greater than zero"
	ErrInvalidIssuerAccount             = "invalid issuer account, %s"
)
