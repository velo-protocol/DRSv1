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

	ErrToGetDataFromDatabase = "can't get data from database"
	ErrToSaveDatabase        = "can't save to database"
	ErrToUpdateDatabase      = "can't update to database"
	ErrToBeginTransaction    = "can't start database transaction"
	ErrToCommitTransaction   = "can't commit database transaction"
	ErrToDeleteRecord        = "can't delete record"
)
