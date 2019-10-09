package models

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestNewPrice(t *testing.T) {
	encodePrice := func(timestamp string, price string) string {
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s_%s", timestamp, price)))
	}
	t.Run("success", func(t *testing.T) {
		price, err := NewPrice(testhelpers.PublicKey1, encodePrice("0", "15000000"))

		assert.NoError(t, err)
		assert.Equal(t, int64(15000000), price.Value)
		assert.Equal(t, int64(0), price.UnixTimestamp)
		assert.Equal(t, testhelpers.PublicKey1, price.Source)
	})
	t.Run("error, bad price source address format", func(t *testing.T) {
		_, err := NewPrice("BAD_PK", encodePrice("0", "15000000"))
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrToDecodeData, "BAD_PK"))
	})
	t.Run("error, bad base 64 string", func(t *testing.T) {
		_, err := NewPrice(testhelpers.PublicKey1, "BAD_BASE_64")
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrToDecodeData, testhelpers.PublicKey1))
	})
	t.Run("error, fail to parse price data", func(t *testing.T) {
		_, err := NewPrice(testhelpers.PublicKey1, base64.StdEncoding.EncodeToString([]byte("BAD_PRICE_DATA")))
		assert.Contains(t, err.Error(), fmt.Sprintf("fail to parse price data from %s", testhelpers.PublicKey1))
	})
	t.Run("error, fail to parse price timestamp", func(t *testing.T) {
		_, err := NewPrice(testhelpers.PublicKey1, encodePrice("BAD-TIMESTAMP", "15000000"))
		assert.Contains(t, err.Error(), fmt.Sprintf("fail to parse price timestamp from %s", testhelpers.PublicKey1))
	})
	t.Run("error, fail to parse price value", func(t *testing.T) {
		_, err := NewPrice(testhelpers.PublicKey1, encodePrice("0", "BAD-PRICE-VALUE"))
		assert.Contains(t, err.Error(), fmt.Sprintf("fail to parse price value from %s", testhelpers.PublicKey1))
	})
}
