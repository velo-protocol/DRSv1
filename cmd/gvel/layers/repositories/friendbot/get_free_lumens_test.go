package friendbot_test

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/friendbot"
	"testing"
)

func TestFb_GetFreeLumens(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		mockedHttpResp := httpmock.NewStringResponder(200, "")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(FakeFriendBotUrl, FakeStellarAddress), mockedHttpResp)

		fbRepo := friendbot.NewFriendBot(FakeFriendBotUrl)

		err := fbRepo.GetFreeLumens(FakeStellarAddress)

		assert.NoError(t, err)
	})

	t.Run("error - failed to get free lumens", func(t *testing.T) {
		mockedHttpResp := httpmock.NewStringResponder(500, "")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(FakeFriendBotUrl, FakeStellarAddress), mockedHttpResp)

		fbRepo := friendbot.NewFriendBot(FakeFriendBotUrl)

		err := fbRepo.GetFreeLumens(FakeStellarAddress)

		assert.Error(t, err)
	})
}
