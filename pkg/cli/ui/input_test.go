package ui_test

import (
	"testing"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/stretchr/testify/assert"
)

func TestPrompt(t *testing.T) {
	go func() {
		keyboard.SimulateKeyPress("Testing")
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	res, err := ui.Prompt("Test?", "")
	assert.NoError(t, err)
	assert.Equal(t, "Testing", res)

	go func() {
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	res, err = ui.Prompt("Test?", "Testing")
	assert.NoError(t, err)
	assert.Equal(t, "Testing", res)
}
