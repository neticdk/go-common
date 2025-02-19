package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSubCommand(t *testing.T) {
	runner := &NoopRunner[any]{}
	cmd := NewSubCommand("test", runner, nil)
	require.NotNil(t, cmd)
	assert.Equal(t, "test", cmd.cmd.Use)
	assert.Nil(t, cmd.cmd.RunE)
}

func TestSubCommandBuilderMethods(t *testing.T) {
	runner := &NoopRunner[any]{}
	builder := NewSubCommand("test", runner, nil)

	builder.WithShortDesc("short description")
	builder.WithLongDesc("long description")
	builder.WithGroupID("group ID")
	builder.WithExample("example usage")
	builder.WithNoArgs()
	cmd := builder.Build()

	assert.Equal(t, "short description", cmd.Short)
	assert.Equal(t, "long description", cmd.Long)
	assert.Equal(t, "group ID", cmd.GroupID)
	assert.Equal(t, "example usage", cmd.Example)
	assert.NotNil(t, cmd.Args)
	assert.NotNil(t, cmd.RunE)
}

func Test_mkRunE(t *testing.T) {
	type arg struct{}

	t.Run("complete_called", func(t *testing.T) {
		completeCalled := false
		runner := &TestRunner[arg]{
			CompleteFunc: func(a arg) error {
				completeCalled = true
				return nil
			},
			ValidateFunc: func(a arg) error { return nil },
			RunFunc:      func(a arg) error { return nil },
		}
		runE := mkRunE(runner, arg{})
		cmd := &cobra.Command{}
		err := runE(cmd, []string{})
		assert.Nil(t, err)
		assert.True(t, completeCalled)
	})

	t.Run("validate_error", func(t *testing.T) {
		runner := &TestRunner[arg]{
			CompleteFunc: func(a arg) error { return nil },
			ValidateFunc: func(a arg) error { return assert.AnError },
			RunFunc:      func(a arg) error { return nil },
		}
		runE := mkRunE(runner, arg{})
		cmd := &cobra.Command{}
		err := runE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("run_called", func(t *testing.T) {
		runCalled := false
		runner := &TestRunner[arg]{
			CompleteFunc: func(a arg) error { return nil },
			ValidateFunc: func(a arg) error { return nil },
			RunFunc:      func(a arg) error { runCalled = true; return nil },
		}

		runE := mkRunE(runner, arg{})
		cmd := &cobra.Command{}
		err := runE(cmd, []string{})
		assert.Nil(t, err)
		assert.True(t, runCalled)
	})
}
