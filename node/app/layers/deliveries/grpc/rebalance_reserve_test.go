package grpc_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/grpc"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestHandler_RebalanceReserve(t *testing.T) {

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		helper.mockUseCase.EXPECT().RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		rebalanceReserveOutput, err := helper.handler.RebalanceReserve(context.Background(), &grpc.RebalanceReserveRequest{})

		assert.Error(t, err)
		assert.Nil(t, rebalanceReserveOutput)
	})
}
