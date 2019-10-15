package friendbot_test

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/friendbot"
	"testing"
)

func TestFb_GetFreeLumens(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		mockedHttpResp := httpmock.NewStringResponder(200, "")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(URL_FAKE_FRIENDBOT, FAKE_ADDRESS), mockedHttpResp)

		fbRepo := friendbot.NewFriendbot(URL_FAKE_FRIENDBOT)

		err := fbRepo.GetFreeLumens(FAKE_ADDRESS)

		assert.NoError(t, err)
	})

	t.Run("error - failed to get free lumens", func(t *testing.T) {
		mockedHttpResp := httpmock.NewStringResponder(500, "")
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf(URL_FAKE_FRIENDBOT, FAKE_ADDRESS), mockedHttpResp)

		fbRepo := friendbot.NewFriendbot(URL_FAKE_FRIENDBOT)

		err := fbRepo.GetFreeLumens(FAKE_ADDRESS)

		assert.Error(t, err)
	})
}
