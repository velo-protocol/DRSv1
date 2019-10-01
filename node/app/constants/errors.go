package constants

var (
	ErrCreateWhiteList          = "can't create white list"
	ErrUnauthorized             = "unauthorized to perform an action"
	ErrRoleNotFound             = "role specified does not exists"
	ErrBadSignature             = "bad signature"
	ErrSignatureNotFound        = "signature not found"
	ErrUnknownVeloOperationType = "unknown velo operation type"
	ErrFormatMissingOperation   = "operation type %s is missing"

	ErrSignatureNotMatchSourceAccount = "the signature and source account does not match"
	ErrFormatSignerNotHavePermission  = "the signer is not found or does not have sufficient permission to perform %s"
	ErrWhiteListAlreadyWhiteListed    = "the address %s has already been whitelisted for the role %s"

	ErrPriceFeederCurrencyMustNotBlank  = "currency must not be blank for price feeder role"
	ErrCurrencyMustBeBlank              = "currency must be blank for non-price feeder role"
	ErrGetSenderAccount                 = "fail to get tx sender account"
	ErrGetDrsAccount                    = "fail to get data of drs account"
	ErrGetRoleListAccount               = "fail to get role list accounts"
	ErrGetTrustedPartnerListDataAccount = "fail to get data of trusted partner list account"
	ErrDerivedKeyPairFromSeed           = "failed to derived KP from seed key"
	ErrBuildAndSignTransaction          = "failed to build and sign tx"
	ErrUnknowRoleSpecified              = "unknown role specified"
	ErrAssetCodeAlreadyBeenUsed         = "asset code %s has already been used"

	ErrCreateTrustedPartnerMetaKeyPair = "failed to create trusted partner meta KP"
	ErrCreateIssuerKeyPair             = "failed to create issuer KP"
	ErrCreateDistributorKeyPair        = "failed to create distributor KP"
)
