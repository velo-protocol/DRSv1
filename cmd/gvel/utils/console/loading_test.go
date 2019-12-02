package console

import (
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestStartLoading(t *testing.T) {
	DefaultLoadWriter = ioutil.Discard
	loadSpinner = nil
	defer func() {
		loadSpinner = nil
		DefaultLoadWriter = color.Output
	}()

	assert.NotPanics(t, func() {
		StartLoading("Now Loading ...")
	})

	assert.NotNil(t, loadSpinner)
	spinner1 := loadSpinner

	// It's fine to start loading twice
	assert.NotPanics(t, func() {
		StartLoading("Now Loading 2 ...")
	})
	assert.NotNil(t, loadSpinner)
	spinner2 := loadSpinner

	assert.False(t, spinner1 == spinner2)
}

func TestStopLoading(t *testing.T) {
	DefaultLoadWriter = ioutil.Discard
	loadSpinner = nil
	defer func() {
		loadSpinner = nil
		DefaultLoadWriter = color.Output
	}()

	StartLoading("Now Loading ...")
	assert.NotPanics(t, func() {
		StopLoading()
	})
	assert.Nil(t, loadSpinner)
}
