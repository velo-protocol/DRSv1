package usecases

import (
	"context"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
)

func (useCase *useCase) UpdatePrice(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError) {
	return nil, nil
	//if err := veloTx.VeloOp.Validate(); err != nil {
	//	return nerrors.ErrInvalidArgument{Message: err.Error()}
	//}
	//
	//txSenderPublicKey := veloTx.TxEnvelope().VeloTx.SourceAccount.Address()
	//txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	//if err != nil {
	//	return nerrors.ErrInvalidArgument{Message: err.Error()}
	//}
	//if veloTx.TxEnvelope().Signatures == nil {
	//	return nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotFound}
	//}
	//if txSenderKeyPair.Hint() != veloTx.TxEnvelope().Signatures[0].Hint {
	//	return nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotMatchSourceAccount}
	//}
	//
	//priceFeederEntity, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhiteListFilter{
	//	StellarPublicKey: pointer.ToString(txSenderPublicKey),
	//	RoleCode:         pointer.ToString(string(vxdr.RolePriceFeeder)),
	//})
	//if err != nil {
	//	return nerrors.ErrInternal{Message: err.Error()}
	//}
	//if priceFeederEntity == nil {
	//	return nerrors.ErrPermissionDenied{
	//		Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpPriceUpdate),
	//	}
	//}
	//
	//createPriceEntity := &entities.CreatePriceEntry{
	//	FeederPublicKey:             txSenderPublicKey,
	//	Asset:                       veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.Asset,
	//	PriceInCurrencyPerAssetUnit: decimal.New(int64(veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.PriceInCurrencyPerAssetUnit), -7),
	//	Currency:                    string(veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.Currency),
	//}
	//
	//_, err = useCase.PriceRepo.CreatePriceEntry(createPriceEntity)
	//if err != nil {
	//	return nerrors.ErrInternal{Message: err.Error()}
	//}
	//
	//return nil
}
