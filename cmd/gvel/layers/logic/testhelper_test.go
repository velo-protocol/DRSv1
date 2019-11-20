package logic_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/logic"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/mocks"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/crypto"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/mocks"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"os"
	"testing"
)

type helper struct {
	logic             logic.Logic
	mockDB            *mocks.MockDbRepo
	mockStellar       *mocks.MockStellarRepo
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
	mockFriendBot := mocks.NewMockStellarRepo(mockCtrl)
	mockVelo := mocks.NewMockVeloRepo(mockCtrl)
	mockVeloClient := mocks.NewMockVeloClient(mockCtrl)
	mockConfiguration := mockutils.NewMockConfiguration(mockCtrl)

	kp, _ := vconvert.SecretKeyToKeyPair("SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX")

	logger, hook := test.NewNullLogger()
	console.Logger = logger

	// to omit what loader print
	console.DefaultLoadWriter = console.Logger.Out

	return helper{
		logic:             logic.NewLogic(mockDB, mockFriendBot, mockVelo, mockConfiguration),
		mockDB:            mockDB,
		mockStellar:       mockFriendBot,
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
