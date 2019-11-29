package vxdr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRole_IsValid(t *testing.T) {
	assert.True(t, RoleTrustedPartner.IsValid())
	assert.True(t, RolePriceFeeder.IsValid())
	assert.True(t, RoleRegulator.IsValid())
	assert.False(t, Role("ADMIN").IsValid())
}
