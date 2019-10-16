package entities

import "github.com/stellar/go/protocols/horizon"

type GetTrustedPartnerFromIssuerAccountInput struct {
	IssuerAccount *horizon.Account
}

type GetTrustedPartnerFromIssuerAccountOutput struct {
	TrustedPartnerAccount *horizon.Account
}
