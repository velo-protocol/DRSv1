package vxdr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsset_IsValid(t *testing.T) {
	assert.True(t, AssetVELO.IsValid())
	assert.False(t, Asset("NOTVELO").IsValid())
}
