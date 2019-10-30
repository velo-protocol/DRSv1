package logic_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
	"github.com/stellar/go/keypair"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/crypto"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/mocks"
	"gitlab.com/velo-labs/cen/libs/convert"
	"os"
	"testing"
)

type helper struct {
	logic             logic.Logic
	mockDB            *mocks.MockDbRepo
	mockFriendBot     *mocks.MockFriendBotRepo
	mockVelo          *mocks.MockVeloRepo
	mockVeloClient    *mocks.MockVeloClient
	mockConfiguration *mockutils.MockConfiguration
	mockController    *gomock.Controller
	keyPair           *keypair.Full
	logHook           *test.Hook
	done              func()
}

func initTest(t *testing.T) helper {
	mockCtrl := gomock.NewController(t)

	mockDB := mocks.NewMockDbRepo(mockCtrl)
	mockFriendBot := mocks.NewMockFriendBotRepo(mockCtrl)
	mockVelo := mocks.NewMockVeloRepo(mockCtrl)
	mockVeloClient := mocks.NewMockVeloClient(mockCtrl)
	mockConfiguration := mockutils.NewMockConfiguration(mockCtrl)

	kp, _ := vconvert.SecretKeyToKeyPair("SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX")

	logger, hook := test.NewNullLogger()
	console.Logger = logger

	return helper{
		logic:             logic.NewLogic(mockDB, mockFriendBot, mockVelo, mockConfiguration),
		mockDB:            mockDB,
		mockFriendBot:     mockFriendBot,
		mockVelo:          mockVelo,
		mockVeloClient:    mockVeloClient,
		mockController:    mockCtrl,
		logHook:           hook,
		mockConfiguration: mockConfiguration,
		keyPair:           kp,
		done: func() {
			viper.Reset()
			_ = os.RemoveAll("./.velo")
		},
	}
}

func arrayOfStellarAccountsBytes() [][]byte {
	return [][]byte{
		stellarAccountsBytes(),
	}
}

func stellarAccountsBytes() []byte {
	stellarAccountBytes, _ := json.Marshal(stellarAccountEntity())
	return stellarAccountBytes
}

func stellarAccountEntity() entity.StellarAccount {
	encryptedSeed, nonce, _ := crypto.Encrypt([]byte("SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"), "password")

	// the passphrase is `password`
	return entity.StellarAccount{
		Address:       "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73",
		EncryptedSeed: encryptedSeed,
		Nonce:         nonce,
	}
}
