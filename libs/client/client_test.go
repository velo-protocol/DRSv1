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
		assert.Contains(t, err.Error(), "cannot connect to Velo Node")
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

func TestGrpc_GetExchangeRate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assetCode := "vTHB"
		assetIssuer := "GCNMY2YGZZNUDMHB3EA36FYWW63ZRAWJX5RQZTZXDLRWCK73H77F264J"
		RedeemableCollateral := "VELO"
		RedeemablePricePerUnit := "1.5"
		helper := initTest(t)

		helper.mockVeloNodeClient.EXPECT().
			GetExchangeRate(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.GetExchangeRateRequest{
				AssetCode: assetCode,
				Issuer:    assetIssuer,
			})).
			Return(&cenGrpc.GetExchangeRateReply{
				AssetCode:              assetCode,
				Issuer:                 assetIssuer,
				RedeemableCollateral:   RedeemableCollateral,
				RedeemablePricePerUnit: RedeemablePricePerUnit,
			}, nil)

		getExchangeRate, err := helper.client.GetExchangeRate(context.Background(), &cenGrpc.GetExchangeRateRequest{
			AssetCode: assetCode,
			Issuer:    assetIssuer,
		})

		assert.NoError(t, err)
		assert.NotNil(t, getExchangeRate)
		assert.Equal(t, assetCode, getExchangeRate.AssetCode)
		assert.Equal(t, assetIssuer, getExchangeRate.Issuer)
		assert.Equal(t, RedeemableCollateral, getExchangeRate.RedeemableCollateral)
		assert.Equal(t, RedeemablePricePerUnit, getExchangeRate.RedeemablePricePerUnit)
	})
	t.Run("error, fail from sdk", func(t *testing.T) {
		helper := initTest(t)

		helper.mockVeloNodeClient.EXPECT().
			GetExchangeRate(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.GetExchangeRateRequest{})).
			Return(nil, errors.New("some error has occurs"))

		getExchangeRate, err := helper.client.GetExchangeRate(context.Background(), &cenGrpc.GetExchangeRateRequest{})

		assert.Error(t, err)
		assert.Nil(t, getExchangeRate)
		assert.Equal(t, "some error has occurs", err.Error())
	})
}

func TestGrpc_GetCollateralHealthCheck(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assetCode := "VELO"
		assetIssuer := "GCNMY2YGZZNUDMHB3EA36FYWW63ZRAWJX5RQZTZXDLRWCK73H77F264J"
		requiredAmount := "350"
		poolAmount := "250"
		helper := initTest(t)

		helper.mockVeloNodeClient.EXPECT().
			GetCollateralHealthCheck(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.GetCollateralHealthCheckRequest{})).
			Return(&cenGrpc.GetCollateralHealthCheckReply{
				AssetCode:      assetCode,
				AssetIssuer:    assetIssuer,
				RequiredAmount: requiredAmount,
				PoolAmount:     poolAmount,
			}, nil)

		getCollateralHealthCheck, err := helper.client.GetCollateralHealthCheck(context.Background(), &cenGrpc.GetCollateralHealthCheckRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, getCollateralHealthCheck)
		assert.Equal(t, assetCode, getCollateralHealthCheck.AssetCode)
		assert.Equal(t, assetIssuer, getCollateralHealthCheck.AssetIssuer)
		assert.Equal(t, requiredAmount, getCollateralHealthCheck.RequiredAmount)
		assert.Equal(t, poolAmount, getCollateralHealthCheck.PoolAmount)
	})
	t.Run("error, fail from sdk", func(t *testing.T) {
		helper := initTest(t)

		helper.mockVeloNodeClient.EXPECT().
			GetCollateralHealthCheck(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.GetCollateralHealthCheckRequest{})).
			Return(nil, errors.New("some error has occurs"))

		getExchangeRate, err := helper.client.GetCollateralHealthCheck(context.Background(), &cenGrpc.GetCollateralHealthCheckRequest{})

		assert.Error(t, err)
		assert.Nil(t, getExchangeRate)
		assert.Equal(t, "some error has occurs", err.Error())
	})
}

func TestGrpc_RebalanceReserve(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assetCode := "VELO"
		assetIssuer := "GCNMY2YGZZNUDMHB3EA36FYWW63ZRAWJX5RQZTZXDLRWCK73H77F264J"
		requiredAmount := "350"
		poolAmount := "250"

		rebalanceCollateral := []*cenGrpc.RebalanceCollateral{{
			AssetCode:      assetCode,
			AssetIssuer:    assetIssuer,
			RequiredAmount: requiredAmount,
			PoolAmount:     poolAmount,
		}}

		helper := initTest(t)

		helper.mockVeloNodeClient.EXPECT().
			RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.RebalanceReserveRequest{})).
			Return(&cenGrpc.RebalanceReserveReply{
				Items:              rebalanceCollateral,
				SignedStellarTxXdr: "AAAAAOOa4D2CPULHjY8jeGzx6g/FL0QUeIpm5juox5lt04wpAAAAZAAFkFAAAAABAAAAAAAAAAEAAAAKMzIzMjA5NjQ2NQAAAAAAAQAAAAEAAAAA45rgPYI9QseNjyN4bPHqD8UvRBR4imbmO6jHmW3TjCkAAAABAAAAAE3j7m7lhZ39noA3ToXWDjJ9QuMmmp/1UaIg0chYzRSlAAAAAAAAAAJMTD+AAAAAAAAAAAFt04wpAAAAQAxFRWcepbQoisfiZ0PG7XhPIBl2ssiD9ymMVpsDyLoHyWXboJLaqibNbiPUHk/KEToTVg7G/JCZ06Mfj0daVAc=",
			}, nil)

		rebalanceReserve, err := helper.client.RebalanceReserve(context.Background(), &cenGrpc.RebalanceReserveRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, rebalanceReserve)
		assert.Equal(t, assetCode, rebalanceCollateral[0].AssetCode)
		assert.Equal(t, assetIssuer, rebalanceCollateral[0].AssetIssuer)
		assert.Equal(t, requiredAmount, rebalanceCollateral[0].RequiredAmount)
		assert.Equal(t, poolAmount, rebalanceCollateral[0].PoolAmount)
	})

	t.Run("error, fail from sdk", func(t *testing.T) {
		helper := initTest(t)

		helper.mockVeloNodeClient.EXPECT().
			RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(&cenGrpc.RebalanceReserveRequest{})).
			Return(nil, errors.New("some error has occurs"))

		rebalanceReserve, err := helper.client.RebalanceReserve(context.Background(), &cenGrpc.RebalanceReserveRequest{})

		assert.Error(t, err)
		assert.Nil(t, rebalanceReserve)
		assert.Equal(t, "some error has occurs", err.Error())
	})
}
