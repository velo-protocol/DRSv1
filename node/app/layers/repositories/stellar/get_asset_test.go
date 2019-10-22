package stellar_test

import (
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/base"
	"github.com/stellar/go/support/render/hal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestRepo_GetAsset(t *testing.T) {
	testhelpers.InitEnv()

	var (
		linkSelf    = "https://horizon-testnet.stellar.org/assets?asset_code=VELO&asset_issuer=GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U&cursor=&limit=200&order=asc"
		linkNext    = "https://horizon-testnet.stellar.org/assets?asset_code=VELO&asset_issuer=GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U&cursor=&limit=200&order=asc"
		linkPrev    = "https://horizon-testnet.stellar.org/assets?asset_code=VELO&asset_issuer=GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U&cursor=&limit=200&order=desc"
		assetType   = "credit_alphanum4"
		pagingToken = "VELO_GB4KSJ74UCNSOCRYLD2HXYNJSX5YQF2X4GDGBNCOIJ3CU76F3TAUOOZJ_credit_alphanum4"
		amount      = "100.0000000"
		order       = horizonclient.Order("asc")
		cursor      = "1"
		limit       = 200 // maximum limit
	)

	mockedAssetStat := []horizon.AssetStat{
		{
			Links: struct {
				Toml hal.Link `json:"toml"`
			}{Toml: hal.Link{
				Href:      "",
				Templated: false,
			}},
			Asset: base.Asset{
				Type:   assetType,
				Code:   string(vxdr.AssetVELO),
				Issuer: testhelpers.TrustedPartnerListPublicKey,
			},
			PT:          pagingToken,
			Amount:      amount,
			NumAccounts: 1,
			Flags: horizon.AccountFlags{
				AuthRequired:  false,
				AuthRevocable: false,
				AuthImmutable: false,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("Assets", horizonclient.AssetRequest{
				ForAssetCode:   string(vxdr.AssetVELO),
				ForAssetIssuer: testhelpers.TrustedPartnerListPublicKey,
				Order:          order,
				Cursor:         cursor,
				Limit:          uint(limit),
			}).
			Return(horizon.AssetsPage{
				Links: hal.Links{
					Self: hal.Link{Href: linkSelf},
					Next: hal.Link{Href: linkNext},
					Prev: hal.Link{Href: linkPrev},
				},
				Embedded: struct {
					Records []horizon.AssetStat
				}{
					Records: mockedAssetStat,
				},
			}, nil)

		asset, err := helper.repo.GetAsset(entities.GetAssetInput{
			AssetCode:   string(vxdr.AssetVELO),
			AssetIssuer: testhelpers.TrustedPartnerListPublicKey,
			Cursor:      pointer.ToString(cursor),
			Order:       pointer.ToString("asc"),
			Limit:       pointer.ToUint(200),
		})

		assert.NoError(t, err)
		assert.NotNil(t, asset)
		assert.Equal(t, pagingToken, asset.Embedded.Records[0].PagingToken())
		assert.Equal(t, string(vxdr.AssetVELO), asset.Embedded.Records[0].Code)
		assert.Equal(t, testhelpers.TrustedPartnerListPublicKey, asset.Embedded.Records[0].Issuer)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "Assets", 1)
	})

	t.Run("success, with no asset", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("Assets", horizonclient.AssetRequest{
				ForAssetCode:   string(vxdr.AssetVELO),
				ForAssetIssuer: testhelpers.TrustedPartnerListPublicKey,
				Order:          order,
				Cursor:         cursor,
				Limit:          uint(limit),
			}).
			Return(horizon.AssetsPage{
				Links: hal.Links{
					Self: hal.Link{Href: linkSelf},
					Next: hal.Link{Href: linkNext},
					Prev: hal.Link{Href: linkPrev},
				},
				Embedded: struct {
					Records []horizon.AssetStat
				}{},
			}, nil)

		asset, err := helper.repo.GetAsset(entities.GetAssetInput{
			AssetCode:   string(vxdr.AssetVELO),
			AssetIssuer: testhelpers.TrustedPartnerListPublicKey,
			Cursor:      pointer.ToString(cursor),
			Order:       pointer.ToString("asc"),
			Limit:       pointer.ToUint(200),
		})

		isNoRecord := []horizon.AssetStat([]horizon.AssetStat(nil))

		assert.NoError(t, err)
		assert.NotNil(t, asset)
		assert.Equal(t, isNoRecord, asset.Embedded.Records)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "Assets", 1)
	})

	t.Run("success, with none option field", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("Assets", horizonclient.AssetRequest{
				ForAssetCode:   string(vxdr.AssetVELO),
				ForAssetIssuer: testhelpers.TrustedPartnerListPublicKey,
			}).
			Return(horizon.AssetsPage{
				Links: hal.Links{
					Self: hal.Link{Href: linkSelf},
					Next: hal.Link{Href: linkNext},
					Prev: hal.Link{Href: linkPrev},
				},
				Embedded: struct {
					Records []horizon.AssetStat
				}{},
			}, nil)

		asset, err := helper.repo.GetAsset(entities.GetAssetInput{
			AssetCode:   string(vxdr.AssetVELO),
			AssetIssuer: testhelpers.TrustedPartnerListPublicKey,
		})

		isNoRecord := []horizon.AssetStat([]horizon.AssetStat(nil))

		assert.NoError(t, err)
		assert.NotNil(t, asset)
		assert.Equal(t, isNoRecord, asset.Embedded.Records)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "Assets", 1)
	})

	t.Run("error, fail to get asset", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("Assets", horizonclient.AssetRequest{
				ForAssetCode:   string(vxdr.AssetVELO),
				ForAssetIssuer: testhelpers.TrustedPartnerListPublicKey,
				Order:          order,
				Cursor:         cursor,
				Limit:          uint(limit),
			}).
			Return(horizon.AssetsPage{}, errors.New("some error has occurs"))

		asset, err := helper.repo.GetAsset(entities.GetAssetInput{
			AssetCode:   string(vxdr.AssetVELO),
			AssetIssuer: testhelpers.TrustedPartnerListPublicKey,
			Cursor:      pointer.ToString(cursor),
			Order:       pointer.ToString("asc"),
			Limit:       pointer.ToUint(200),
		})

		assert.Error(t, err)
		assert.Nil(t, asset)
		assert.Equal(t, "fail to get asset VELO: some error has occurs", err.Error())

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "Assets", 1)
	})
}
