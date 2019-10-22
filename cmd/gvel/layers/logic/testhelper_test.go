package logic_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/mocks"
	"testing"
)

type helper struct {
	logic          logic.Logic
	mockDB         *mocks.MockDbRepo
	mockFriendBot  *mocks.MockFriendBotRepo
	mockController *gomock.Controller
}

func initTest(t *testing.T) helper {
	mockCtrl := gomock.NewController(t)

	mockDB := mocks.NewMockDbRepo(mockCtrl)
	mockFriendBot := mocks.NewMockFriendBotRepo(mockCtrl)

	return helper{
		logic:          logic.NewLogic(mockDB, mockFriendBot),
		mockDB:         mockDB,
		mockFriendBot:  mockFriendBot,
		mockController: mockCtrl,
	}
}

func stellarAccountsBytes() [][]byte {
	stellarAccountBytes, _ := json.Marshal(entity.StellarAccount{
		EncryptedSeed: []byte("fake-seed"),
		Nonce:         []byte("aaaa"),
	})

	return [][]byte{
		stellarAccountBytes,
	}
}
