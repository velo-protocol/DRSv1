package vxdr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrency_IsValid(t *testing.T) {
	assert.True(t, CurrencySGD.IsValid())
	assert.True(t, CurrencyUSD.IsValid())
	assert.True(t, CurrencyTHB.IsValid())
	assert.False(t, Currency("JPY").IsValid())
}
