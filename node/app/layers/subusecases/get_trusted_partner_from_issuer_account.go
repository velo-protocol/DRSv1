package subusecases

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
)

func (subUseCase *subUseCase) GetTrustedPartnerFromIssuerAccount(ctx context.Context, input *entities.GetTrustedPartnerFromIssuerAccountInput) (*entities.GetTrustedPartnerFromIssuerAccountOutput, error) {
	// derive trusted partner address
	var trustedPartnerAddress string
	var drsKeyFound, issuerKeyFound bool // make sure the other two keys are what we expected
	for _, signer := range input.IssuerAccount.Signers {
		if signer.Key != env.DrsPublicKey && signer.Key != input.IssuerAccount.AccountID {
			trustedPartnerAddress = signer.Key
		}

		if signer.Key == env.DrsPublicKey && signer.Weight == 1 {
			drsKeyFound = true
		}
		if signer.Key == input.IssuerAccount.AccountID && signer.Weight == 0 {
			issuerKeyFound = true
		}
	}
	if trustedPartnerAddress == "" || !drsKeyFound || !issuerKeyFound {
		return nil, errors.Errorf(constants.ErrInvalidIssuerAccount, "no drs as signer")
	}

	// get trusted partner account
	trustedPartnerAccount, err := subUseCase.StellarRepo.GetAccount(trustedPartnerAddress)
	if err != nil {
		return nil, err
	}

	return &entities.GetTrustedPartnerFromIssuerAccountOutput{
		TrustedPartnerAccount: trustedPartnerAccount,
	}, nil
}
