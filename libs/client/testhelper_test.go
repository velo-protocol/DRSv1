package vclient

import (
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	mock_grpc "gitlab.com/velo-labs/cen/grpc/mocks"
	"gitlab.com/velo-labs/cen/libs/convert"
	"google.golang.org/grpc"
	"testing"
)

const (
	clientPublicKey       = "GC3YGJSKM67UGPCQ37R7CREYONJJG7WQAVMIQJQ7QIF2C2IQVVVW3NJP"
	clientSecretKey       = "SDRBFTQ2K72EG2O54WQGFZCW6OLGLPNIQIGJJYGV5CGKA6KBHDSSNIAF"
	whitelistingPublicKey = "GADZ2M2F5ENOIU2GI7QRCVX4UQNR4KVLYFRXHGI6DQBQLYK4TNZ33ITQ"
	whitelistingSecretKey = "SDINNNFBMSQQK5OW7AXSNNCURC6G32AKQNO35CHJLWVKXYSQHCNZPR4I"
	drsPublicKey          = "GCWX4N3SKRBNLIRL5TZCJZGDCUYG4LGZSBED7HSHFC3HPXA6DR5UJN4A"
	drsSecretKey          = "SBNQIMGPERFYDHQ54NLUVIT74DDHEGG6QUG5SKQ3J24KMSRXBTTYZN2Z"
)

var (
	clientKp, _       = vconvert.SecretKeyToKeyPair(clientSecretKey)
	whitelistingKp, _ = vconvert.SecretKeyToKeyPair(whitelistingSecretKey)
	drsKp, _          = vconvert.SecretKeyToKeyPair(drsSecretKey)

	getSimpleBumpTx = func() txnbuild.Transaction {
		return txnbuild.Transaction{
			SourceAccount: &horizon.Account{
				AccountID: clientPublicKey,
				Sequence:  "100",
			},
			Operations: []txnbuild.Operation{
				&txnbuild.BumpSequence{
					BumpTo: 2,
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: whitelistingPublicKey,
					},
				},
			},
			Network:    "Test",
			Timebounds: txnbuild.NewInfiniteTimeout(),
		}
	}
	getSimpleBumpTxXdr = func(kps ...*keypair.Full) string {
		tx := getSimpleBumpTx()
		xdr, _ := tx.BuildSignEncode(kps...)
		return xdr
	}
)

type helper struct {
	mockHorizonClient  *horizonclient.MockClient
	mockVeloNodeClient *mock_grpc.MockVeloNodeClient
	client             *Client
}

func initTest(t *testing.T) *helper {
	mockHorizonClient := new(horizonclient.MockClient)

	mockVeloNodeClientController := gomock.NewController(t)
	mockVeloNodeClient := mock_grpc.NewMockVeloNodeClient(mockVeloNodeClientController)

	return &helper{
		mockHorizonClient:  mockHorizonClient,
		mockVeloNodeClient: mockVeloNodeClient,
		client: &Client{
			horizonClient:     mockHorizonClient,
			networkPassphrase: "Test",
			keyPair:           clientKp,
			veloNodeClient:    mockVeloNodeClient,
			grpcConnection:    new(grpc.ClientConn),
		},
	}
}
