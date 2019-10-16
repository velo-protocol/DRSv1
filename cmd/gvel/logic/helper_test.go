package logic_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	mocksDB "gitlab.com/velo-labs/cen/cmd/gvel/db/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	mocksFriendbot "gitlab.com/velo-labs/cen/cmd/gvel/friendbot/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/logic"
	"testing"
)

type helper struct {
	logic          logic.Logic
	mockDB         *mocksDB.MockDB
	mockFriendbot  *mocksFriendbot.MockRepository
	mockController *gomock.Controller
}

func initTest(t *testing.T) helper {
	mockCtrl := gomock.NewController(t)

	mockDB := mocksDB.NewMockDB(mockCtrl)
	mockFriendbot := mocksFriendbot.NewMockRepository(mockCtrl)

	return helper{
		logic:          logic.NewLogic(mockDB, mockFriendbot),
		mockDB:         mockDB,
		mockFriendbot:  mockFriendbot,
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
