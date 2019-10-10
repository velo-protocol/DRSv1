package vclient

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/support/render/problem"
	"github.com/stretchr/testify/assert"
	cenGrpc "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestNewPublicClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, err := NewPublicClient(&grpc.ClientConn{}, clientSecretKey)
		assert.NoError(t, err)
		assert.Equal(t, horizonclient.DefaultPublicNetClient, client.horizonClient)
		assert.Equal(t, network.PublicNetworkPassphrase, client.networkPassphrase)
		assert.Equal(t, clientKp, client.keyPair)
	})
	t.Run("error, fail to key pair from secret key", func(t *testing.T) {
		_, err := NewPublicClient(&grpc.ClientConn{}, "BAD_SECRET_KEY")
		assert.Error(t, err)
	})
}
func TestNewTestNetClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, err := NewTestNetClient(&grpc.ClientConn{}, clientSecretKey)
		assert.NoError(t, err)
		assert.Equal(t, horizonclient.DefaultTestNetClient, client.horizonClient)
		assert.Equal(t, network.TestNetworkPassphrase, client.networkPassphrase)
		assert.Equal(t, clientKp, client.keyPair)
	})
	t.Run("error, fail to key pair from secret key", func(t *testing.T) {
		_, err := NewTestNetClient(&grpc.ClientConn{}, "BAD_SECRET_KEY")
		assert.Error(t, err)
	})
}

func TestClient_executeVeloTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockVeloNodeClient.EXPECT().
			SubmitVeloTx(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.VeloTxRequest{})).
			Return(&cenGrpc.VeloTxReply{
				SignedStellarTxXdr: getSimpleBumpTxXdr(drsKp),
				Message:            "Success",
			}, nil)
		helper.mockHorizonClient.
			On("SubmitTransactionXDR", getSimpleBumpTxXdr(drsKp, clientKp)).
			Return(horizon.TransactionSuccess{
				Result: "AAAA...",
			}, nil)

		result, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{
			Address: whitelistingPublicKey,
			Role:    string(vxdr.RoleRegulator),
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, result.Result)
	})
	t.Run("error, fail to build, sign or encode velo tx", func(t *testing.T) {
		helper := initTest(t)
		_, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{})

		assert.Error(t, err)
	})
	t.Run("error, cannot connect to VeloCen via gRPC", func(t *testing.T) {
		helper := initTest(t)
		helper.mockVeloNodeClient.EXPECT().
			SubmitVeloTx(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.VeloTxRequest{})).
			Return(nil, status.Error(codes.Unavailable, "some error has occurred"))

		_, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{
			Address: whitelistingPublicKey,
			Role:    string(vxdr.RoleRegulator),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot connect to VeloCen via gRPC")
	})
	t.Run("error, velo node client returns an error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockVeloNodeClient.EXPECT().
			SubmitVeloTx(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.VeloTxRequest{})).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{
			Address: whitelistingPublicKey,
			Role:    string(vxdr.RoleRegulator),
		})

		assert.Error(t, err)
	})
	t.Run("error, fail to parse signed velo tx xdr", func(t *testing.T) {
		helper := initTest(t)
		helper.mockVeloNodeClient.EXPECT().
			SubmitVeloTx(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.VeloTxRequest{})).
			Return(&cenGrpc.VeloTxReply{
				SignedStellarTxXdr: "AAAA...BAD_XDR",
				Message:            "Success",
			}, nil)

		_, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{
			Address: whitelistingPublicKey,
			Role:    string(vxdr.RoleRegulator),
		})

		assert.Error(t, err)
	})
	t.Run("error, horizon response with an error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockVeloNodeClient.EXPECT().
			SubmitVeloTx(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.VeloTxRequest{})).
			Return(&cenGrpc.VeloTxReply{
				SignedStellarTxXdr: getSimpleBumpTxXdr(drsKp),
				Message:            "Success",
			}, nil)
		helper.mockHorizonClient.
			On("SubmitTransactionXDR", getSimpleBumpTxXdr(drsKp, clientKp)).
			Return(horizon.TransactionSuccess{}, &horizonclient.Error{
				Problem: problem.BadRequest,
			})

		_, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{
			Address: whitelistingPublicKey,
			Role:    string(vxdr.RoleRegulator),
		})

		assert.Error(t, err)
		assert.IsType(t, &horizonclient.Error{}, err)
	})
	t.Run("error, cannot connect to horizon", func(t *testing.T) {
		helper := initTest(t)
		helper.mockVeloNodeClient.EXPECT().
			SubmitVeloTx(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.VeloTxRequest{})).
			Return(&cenGrpc.VeloTxReply{
				SignedStellarTxXdr: getSimpleBumpTxXdr(drsKp),
				Message:            "Success",
			}, nil)
		helper.mockHorizonClient.
			On("SubmitTransactionXDR", getSimpleBumpTxXdr(drsKp, clientKp)).
			Return(horizon.TransactionSuccess{}, errors.New("some error has occurred"))

		_, err := helper.client.executeVeloTx(context.Background(), &vtxnbuild.Whitelist{
			Address: whitelistingPublicKey,
			Role:    string(vxdr.RoleRegulator),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot connect to horizon")
	})
}
