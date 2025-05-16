package cmd

import (
	"context"
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
			SetupFlagsFunc: func(ctx context.Context, cmd *cobra.Command) error { return nil },
			CompleteFunc: func(ctx context.Context, a arg) error {
				completeCalled = true
				return nil
			},
			ValidateFunc: func(ctx context.Context, a arg) error { return nil },
			RunFunc:      func(ctx context.Context, a arg) error { return nil },
		}
		runE := mkRunE(runner, arg{})
		cmd := &cobra.Command{}
		err := runE(cmd, []string{})
		assert.Nil(t, err)
		assert.True(t, completeCalled)
	})

	t.Run("validate_error", func(t *testing.T) {
		runner := &TestRunner[arg]{
			SetupFlagsFunc: func(ctx context.Context, cmd *cobra.Command) error { return nil },
			CompleteFunc:   func(ctx context.Context, a arg) error { return nil },
			ValidateFunc:   func(ctx context.Context, a arg) error { return assert.AnError },
			RunFunc:        func(ctx context.Context, a arg) error { return nil },
		}
		runE := mkRunE(runner, arg{})
		cmd := &cobra.Command{}
		err := runE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("run_called", func(t *testing.T) {
		runCalled := false
		runner := &TestRunner[arg]{
			SetupFlagsFunc: func(ctx context.Context, cmd *cobra.Command) error { return nil },
			CompleteFunc:   func(ctx context.Context, a arg) error { return nil },
			ValidateFunc:   func(ctx context.Context, a arg) error { return nil },
			RunFunc:        func(ctx context.Context, a arg) error { runCalled = true; return nil },
		}

		runE := mkRunE(runner, arg{})
		cmd := &cobra.Command{}
		err := runE(cmd, []string{})
		assert.Nil(t, err)
		assert.True(t, runCalled)
	})
}

func TestArgsHelpers(t *testing.T) {
	t.Run("NoArgsReturnsErrorWhenArgsProvided", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithNoArgs()
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1"})
		assert.Error(t, err)
	})

	t.Run("NoArgsSucceedsWhenNoArgsProvided", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithNoArgs()
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("ExactArgsReturnsErrorWhenArgsCountMismatch", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithExactArgs(2)
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1"})
		assert.Error(t, err)
	})

	t.Run("ExactArgsSucceedsWhenArgsCountMatches", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithExactArgs(2)
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1", "arg2"})
		assert.NoError(t, err)
	})

	t.Run("MinArgsReturnsErrorWhenArgsCountLessThanMin", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithMinArgs(2)
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1"})
		assert.Error(t, err)
	})

	t.Run("MinArgsSucceedsWhenArgsCountAtLeastMin", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithMinArgs(2)
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1", "arg2"})
		assert.NoError(t, err)
	})

	t.Run("MaxArgsReturnsErrorWhenArgsCountMoreThanMax", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithMaxArgs(2)
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1", "arg2", "arg3"})
		assert.Error(t, err)
	})

	t.Run("MaxArgsSucceedsWhenArgsCountAtMostMax", func(t *testing.T) {
		builder := NewSubCommand("test", &NoopRunner[any]{}, nil)
		builder.WithMaxArgs(2)
		cmd := builder.Build()

		err := cmd.Args(cmd, []string{"arg1", "arg2"})
		assert.NoError(t, err)
	})
}
