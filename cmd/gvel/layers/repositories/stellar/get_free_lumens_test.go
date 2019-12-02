package stellar_test

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/stellar"
	"testing"
)

func TestFb_GetFreeLumens(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		mockedHttpResp := httpmock.NewStringResponder(200, "")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(FakeFriendBotUrl, FakeStellarAddress), mockedHttpResp)

		stellarRepo := stellar.NewStellarWithClientInterface(FakeHorizonUrl, nil)

		err := stellarRepo.GetFreeLumens(FakeStellarAddress)

		assert.NoError(t, err)
	})

	t.Run("error - failed to get free lumens", func(t *testing.T) {
		mockedHttpResp := httpmock.NewStringResponder(500, "")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(FakeFriendBotUrl, FakeStellarAddress), mockedHttpResp)

		stellarRepo := stellar.NewStellarWithClientInterface(FakeHorizonUrl, nil)

		err := stellarRepo.GetFreeLumens(FakeStellarAddress)

		assert.Error(t, err)
	})
}
