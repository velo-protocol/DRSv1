package grpc_test

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/velo-labs/cen/grpc"
	_handler "gitlab.com/velo-labs/cen/node/app/layers/deliveries/grpc"
	"gitlab.com/velo-labs/cen/node/app/layers/mocks"
	"testing"
)

const (
	publicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
	secretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
	publicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
	//secretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
)

type helper struct {
	handler        grpc.VeloNodeServer
	mockUseCase    mocks.MockUseCase
	mockController *gomock.Controller
}

func initTest(t *testing.T) helper {
	mockCtrl := gomock.NewController(t)
	mockUseCase := mocks.NewMockUseCase(mockCtrl)
	return helper{
		handler:        _handler.InitHandler(mockUseCase),
		mockUseCase:    *mockUseCase,
		mockController: mockCtrl,
	}
}
