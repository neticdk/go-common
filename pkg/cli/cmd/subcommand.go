package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// SubCommandRunner is an interface for a subcommand runner
type SubCommandRunner[T any] interface {
	// Complete performs any setup or completion of arguments
	Complete(ctx context.Context, arg T) error

	// Validate checks if the arguments are valid
	// Returns error if validation fails
	Validate(ctx context.Context, arg T) error

	// Run executes the command with the given arguments
	// Returns error if execution fails
	Run(ctx context.Context, arg T) error
}

// SubCommandBuilder is a builder for a subcommand
type SubCommandBuilder[T any] struct {
	cmd       *cobra.Command
	runner    SubCommandRunner[T]
	runnerArg T
}

// NewSubCommand creates a new subcommand
// args can be used to pass arguments to the runner
func NewSubCommand[T any](
	name string,
	runner SubCommandRunner[T],
	runnerArg T,
) *SubCommandBuilder[T] {
	c := &cobra.Command{
		Use: name,
	}

	return &SubCommandBuilder[T]{
		cmd:       c,
		runner:    runner,
		runnerArg: runnerArg,
	}
}

// Build builds the subcommand
func (b *SubCommandBuilder[T]) Build() *cobra.Command {
	b.cmd.RunE = mkRunE(b.runner, b.runnerArg)
	return b.cmd
}

// WithShortDesc sets the short description of the subcommand
func (b *SubCommandBuilder[T]) WithShortDesc(desc string) *SubCommandBuilder[T] {
	b.cmd.Short = desc
	return b
}

// WithLongDesc sets the long description of the subcommand
func (b *SubCommandBuilder[T]) WithLongDesc(desc string) *SubCommandBuilder[T] {
	b.cmd.Long = desc
	return b
}

// WithGroupID sets the group id of the subcommand
func (b *SubCommandBuilder[T]) WithGroupID(group string) *SubCommandBuilder[T] {
	b.cmd.GroupID = group
	return b
}

// WithExample sets the example of the subcommand
func (b *SubCommandBuilder[T]) WithExample(example string) *SubCommandBuilder[T] {
	b.cmd.Example = example
	return b
}

// WithAliases sets the aliases of the subcommand
func (b *SubCommandBuilder[T]) WithAliases(aliases ...string) *SubCommandBuilder[T] {
	b.cmd.Aliases = aliases
	return b
}

// WithNoArgs causes the subcommand to return an error if any arguments are passed
func (b *SubCommandBuilder[T]) WithNoArgs() *SubCommandBuilder[T] {
	b.cmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf(
				"%q does not accept any arguments\n\nUsage: %s",
				cmd.CommandPath(),
				cmd.UseLine(),
			)
		}
		return nil
	}
	return b
}

// WithExactArgs causes the subcommand to return an error if there are not exactly n arguments
func (b *SubCommandBuilder[T]) WithExactArgs(n int) *SubCommandBuilder[T] {
	b.cmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf(
				"%q requires exactly %d arguments\n\nUsage: %s",
				cmd.CommandPath(),
				n,
				cmd.UseLine(),
			)
		}
		return nil
	}
	return b
}

// WithMinArgs causes the subcommand to return an error if there is not at least n arguments
func (b *SubCommandBuilder[T]) WithMinArgs(n int) *SubCommandBuilder[T] {
	b.cmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf(
				"%q requires at least %d arguments\n\nUsage: %s",
				cmd.CommandPath(),
				n,
				cmd.UseLine(),
			)
		}
		return nil
	}
	return b
}

// WithMaxArgs causes the subcommand to return an error if there are more than n arguments
func (b *SubCommandBuilder[T]) WithMaxArgs(n int) *SubCommandBuilder[T] {
	b.cmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) > n {
			return fmt.Errorf(
				"%q accepts at most %d arguments\n\nUsage: %s",
				cmd.CommandPath(),
				n,
				cmd.UseLine(),
			)
		}
		return nil
	}
	return b
}

// mkRunE creates a RunE function for a subcommand
// runnerArgs are passed to the runner.Complete, runner.Validate and runner.Run methods
// it is up to the runner to decide how to use these arguments
func mkRunE[T any](runner SubCommandRunner[T], runnerArg T) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := runner.Complete(ctx, runnerArg); err != nil {
			return err
		}
		if err := runner.Validate(ctx, runnerArg); err != nil {
			return err
		}
		return runner.Run(ctx, runnerArg)
	}
}

// TestRunner is a test runner for commands
type TestRunner[T any] struct {
	CompleteFunc func(T) error
	ValidateFunc func(T) error
	RunFunc      func(T) error
}

func (tr *TestRunner[T]) Complete(ctx context.Context, runnerArg T) error {
	return tr.CompleteFunc(runnerArg)
}

func (tr *TestRunner[T]) Validate(ctx context.Context, runnerArg T) error {
	return tr.ValidateFunc(runnerArg)
}

func (tr *TestRunner[T]) Run(ctx context.Context, runnerArg T) error {
	return tr.RunFunc(runnerArg)
}

type NoopRunner[T any] struct{}

func (o *NoopRunner[T]) Complete(ctx context.Context, runnerArg T) error { return nil }
func (o *NoopRunner[T]) Validate(ctx context.Context, runnerArg T) error { return nil }
func (o *NoopRunner[T]) Run(ctx context.Context, runnerArg T) error      { return nil }
