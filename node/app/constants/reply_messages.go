package constants

const (
	replyWarningSuffix = "Please sign the transaction and submit to Stellar for the operation to be completed. The transaction will expire in the next 15 minutes."

	ReplyWhitelistSuccess        = "Whitelist operation of address %s as a %s returned. " + replyWarningSuffix
	ReplySetupCreditSuccess      = "Setup credit operation returned. " + replyWarningSuffix
	ReplyPriceUpdateSuccess      = "Price update transaction returned. " + replyWarningSuffix
	ReplyMintCreditSuccess       = "Mint price-stable credit transaction returned. You can mint %s %s in exchange for %s %s. " + replyWarningSuffix
	ReplyRedeemCreditSuccess     = "Redemption transaction returned. " + replyWarningSuffix
	ReplyRebalanceReserveSuccess = "Rebalance reserve transaction returned. " + replyWarningSuffix
)
