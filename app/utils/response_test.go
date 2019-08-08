package utils_test

import (
	"gitlab.com/lightnet-thailand/velo-operators/report-service/app/utils"
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
