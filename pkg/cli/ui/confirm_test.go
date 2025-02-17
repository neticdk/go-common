package ui_test

import (
	"testing"

	"atomicgo.dev/keyboard"
	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/stretchr/testify/assert"
)

func TestConfirm(t *testing.T) {
	go func() {
		keyboard.SimulateKeyPress('Y')
	}()

	res := ui.Confirm("Are you sure?")
	assert.True(t, res)

	go func() {
		keyboard.SimulateKeyPress('N')
	}()

	res = ui.Confirm("Are you sure?")
	assert.False(t, res)
}
