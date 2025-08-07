package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neticdk/go-common/pkg/tui/help"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	GroupBase  = "group-base"
	GroupOther = "group-other"
)

type InitFunc = func(cmd *cobra.Command, args []string) error

type RootCommandBuilder struct {
	cmd *cobra.Command
	ec  *ExecutionContext
}

// NewRootCommand creates a new root command
// Use this function to create a new root command that is used to add subcommands
//
// It supports the following features:
//
// - Adding global flags to the command
// - Automatically reads configuration from the configuration file
// - Automatically reads configuration from environment variables
// - Automatically binds the configuration to the command's flags
// - Automatically sets the log level based on the configuration
// - Automatically sets the log format based on the configuration
//
// It uses ec.AppName as the base name for the configuration file and environment variables
// initFunc is a function that is called before the command is executed
// It can be used to add more context or do other initializations
func NewRootCommand(ec *ExecutionContext) *RootCommandBuilder {
	if ec == nil {
		panic("execution context is required")
	}
	if ec.AppName == "" {
		panic("app name is required")
	}

	c := &cobra.Command{
		Use:                   fmt.Sprintf("%s [command] [flags]", ec.AppName),
		DisableFlagsInUseLine: true,
		Short:                 ec.ShortDescription,
		Long:                  ec.LongDescription,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Version:               ec.Version,
		PersistentPreRunE:     defaultInitFunc(ec),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	// Update default execution context command
	// For subcommands, this is set to the subcommand's command though
	// PersistentPreRunE (defaultInitFunc)
	ec.Command = c

	AddPersistentFlags(c, ec)
	ec.SetLogLevel()

	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		panic(fmt.Errorf("binding flags: %w", err))
	}

	// 💃👄🪄✨️🌈 Add fabulous glamour 🌈✨️🪄👄💃
	if !ec.PFlags.NoColor && ec.IsTerminal {
		help.AddToRootCmd(c)
	}

	c.SetOut(ec.Stdout)
	c.SetErr(ec.Stderr)

	if ec.PFlags.Quiet {
		pterm.DisableOutput()
	}

	c.AddGroup(
		&cobra.Group{
			ID:    GroupBase,
			Title: "Basic Commands:",
		},
		&cobra.Group{
			ID:    GroupOther,
			Title: "Other Commands:",
		},
	)
	c.SetHelpCommandGroupID(GroupOther)
	c.SetCompletionCommandGroupID(GroupOther)

	c.AddCommand(
		GenDocsCommand(ec),
	)

	return &RootCommandBuilder{
		cmd: c,
		ec:  ec,
	}
}

// Build builds the root command
func (b *RootCommandBuilder) Build() *cobra.Command {
	return b.cmd
}

// WithInitFunc adds an init function to the root command
// This sets the PersistentPreRunE function of the command
func (b *RootCommandBuilder) WithInitFunc(fn InitFunc) *RootCommandBuilder {
	initFunc := func(cmd *cobra.Command, args []string) error {
		if err := defaultInitFunc(b.ec)(cmd, args); err != nil {
			return err
		}
		return fn(cmd, args)
	}
	b.cmd.PersistentPreRunE = initFunc
	return b
}

// WithPreRunFunc adds persistent flags to the root command
func (b *RootCommandBuilder) WithPersistentFlags(flags *pflag.FlagSet) *RootCommandBuilder {
	b.cmd.PersistentFlags().AddFlagSet(flags)
	return b
}

// WithExample sets the example of the root command
func (b *RootCommandBuilder) WithExample(example string) *RootCommandBuilder {
	b.cmd.Example = example
	return b
}

func (b *RootCommandBuilder) WithNoSubCommands() *RootCommandBuilder {
	b.cmd.Use = fmt.Sprintf("%s [flags]", b.ec.AppName)
	b.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = cmd.Help()
			return fmt.Errorf("no arguments provided")
		}
		return nil
	}
	return b
}

// initConfig ensures that precedence of configuration setting is correct
// precedence:
// flag -> environment -> configuration file value -> flag default
func initConfig(appName string, cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(appName)
	v.AddConfigPath(".")
	if homeDir, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(homeDir, ".config", appName))
	}
	if dir, err := os.UserConfigDir(); err == nil {
		v.AddConfigPath(filepath.Join(dir, appName))
	}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.SetEnvPrefix(strings.ToUpper(appName))
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	return bindFlags(cmd, v)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			if setErr := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); setErr != nil {
				errs = append(errs, fmt.Errorf("setting flag %q: %w", f.Name, setErr))
			}
		}
	})
	if len(errs) > 0 {
		return fmt.Errorf("binding flags: %v", errs)
	}
	return nil
}

func defaultInitFunc(ec *ExecutionContext) InitFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := initConfig(ec.AppName, cmd); err != nil {
			return err
		}
		ec.Command = cmd
		ec.CommandArgs = args
		ec.SetLogLevel()
		return nil
	}
}
