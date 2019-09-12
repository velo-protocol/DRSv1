package utils_test

import (
	"gitlab.com/velo-labs/cen/app/utils"
	"testing"
)

func TestUtils_NewErrorResponse(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		utils.NewErrorResponse("something bad happened")
	})
}

func TestUtils_NewSuccessResponse(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		utils.NewSuccessResponse(nil)
	})
}
