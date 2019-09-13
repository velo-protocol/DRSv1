package utils_test

import (
	"github.com/stretchr/testify/assert"
	test_helpers "gitlab.com/velo-labs/cen/node/app/test_helpers"
	"gitlab.com/velo-labs/cen/node/app/utils"
	"testing"
)

func TestUtils_KpFromSeedString(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		randStellarAccount := test_helpers.GetRandStellarAccount()

		kp, err := utils.KpFromSeedString(randStellarAccount.Seed())

		assert.NoError(t, err)
		assert.Equal(t, randStellarAccount.Address(), kp.Address())
		assert.Equal(t, randStellarAccount.Seed(), kp.Seed())
	})

	t.Run("error - seed is not in the right format", func(t *testing.T) {
		kp, err := utils.KpFromSeedString("hello world")

		assert.Error(t, err)
		assert.Nil(t, kp)
	})
}
