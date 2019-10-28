package logic_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/mocks"
	"os"
	"testing"
)

type helper struct {
	logic             logic.Logic
	mockDB            *mocks.MockDbRepo
	mockFriendBot     *mocks.MockFriendBotRepo
	mockConfiguration *mockutils.MockConfiguration
	mockController    *gomock.Controller
	done              func()
}

func initTest(t *testing.T) helper {
	mockCtrl := gomock.NewController(t)

	mockDB := mocks.NewMockDbRepo(mockCtrl)
	mockFriendBot := mocks.NewMockFriendBotRepo(mockCtrl)
	mockConfiguration := mockutils.NewMockConfiguration(mockCtrl)

	logger, _ := test.NewNullLogger()
	console.Logger = logger

	return helper{
		logic:             logic.NewLogic(mockDB, mockFriendBot, mockConfiguration),
		mockDB:            mockDB,
		mockFriendBot:     mockFriendBot,
		mockController:    mockCtrl,
		mockConfiguration: mockConfiguration,
		done: func() {
			viper.Reset()
			_ = os.RemoveAll("./.velo")
		},
	}
}

func stellarAccountsBytes() [][]byte {
	stellarAccountBytes, _ := json.Marshal(entity.StellarAccount{
		Address:       "GA...",
		EncryptedSeed: []byte("fake-seed"),
		Nonce:         []byte("aaaa"),
	})

	return [][]byte{
		stellarAccountBytes,
	}
}
