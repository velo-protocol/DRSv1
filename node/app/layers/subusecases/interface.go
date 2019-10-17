package subusecases

import (
	"context"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type SubUseCase interface {
	GetIssuerAccount(ctx context.Context, input *entities.GetIssuerAccountInput) (*entities.GetIssuerAccountOutput, error)
	GetTrustedPartnerFromIssuerAccount(ctx context.Context, input *entities.GetTrustedPartnerFromIssuerAccountInput) (*entities.GetTrustedPartnerFromIssuerAccountOutput, error)
}
