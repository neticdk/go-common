package cmd

import "github.com/spf13/cobra"

// SubCommandRunner is an interface for a subcommand runner
// These commands are the commands run in that order
type SubCommandRunner interface {
	Complete(...any)
	Validate(...any) error
	Run(...any) error
}

// NewSubCommand creates a new subcommand
// args are passed to the runner.Complete, runner.Validate and runner.Run methods
// it is up to the runner to decide how to use these arguments
func NewSubCommand(name, shortDesc, longDesc, group string, runner SubCommandRunner, args ...any) *cobra.Command {
	c := &cobra.Command{
		Use:     name,
		Short:   shortDesc,
		Long:    longDesc,
		GroupID: group,
		RunE:    mkRunE(runner, args),
	}
	return c
}

// mkRunE creates a RunE function for a subcommand
// runnerArgs are passed to the runner.Complete, runner.Validate and runner.Run methods
// it is up to the runner to decide how to use these arguments
func mkRunE(runner SubCommandRunner, runnerArgs ...any) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		runner.Complete(runnerArgs...)
		if err := runner.Validate(runnerArgs...); err != nil {
			return err
		}
		return runner.Run(runnerArgs...)
	}
}
