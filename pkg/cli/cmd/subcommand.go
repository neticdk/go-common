package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// SubCommandRunner is an interface for a subcommand runner
type SubCommandRunner[T any] interface {
	// Complete performs any setup or completion of arguments
	Complete(ctx context.Context, arg T)

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

// Build builds the subcommand
func (b *SubCommandBuilder[T]) Build() *cobra.Command {
	b.cmd.RunE = mkRunE(b.runner, b.runnerArg)
	return b.cmd
}

// mkRunE creates a RunE function for a subcommand
// runnerArgs are passed to the runner.Complete, runner.Validate and runner.Run methods
// it is up to the runner to decide how to use these arguments
func mkRunE[T any](runner SubCommandRunner[T], runnerArg T) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		runner.Complete(ctx, runnerArg)
		if err := runner.Validate(ctx, runnerArg); err != nil {
			return err
		}
		return runner.Run(ctx, runnerArg)
	}
}
