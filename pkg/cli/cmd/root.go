package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neticdk/go-common/pkg/cli/context"
	"github.com/neticdk/go-common/pkg/cli/flags"
	"github.com/neticdk/go-common/pkg/tui/help"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	GroupBase  = "group-base"
	GroupOther = "group-other"
)

type InitFunc = func(cmd *cobra.Command, args []string) error

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
func NewRootCommand(ec *context.ExecutionContext, initFunc InitFunc) *cobra.Command {
	c := &cobra.Command{
		Use:                   fmt.Sprintf("%s [command] [flags]", ec.AppName),
		DisableFlagsInUseLine: true,
		Short:                 ec.ShortDescription,
		Long:                  ec.LongDescription,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Version:               ec.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initConfig(ec.AppName, cmd); err != nil {
				return err
			}
			ec.Command = cmd
			ec.SetLogLevel()

			return initFunc(cmd, args)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// 💃👄🪄✨️🌈 Add fabulous glamour 🌈✨️🪄👄💃
	if !ec.NoColor && ec.IsTerminal {
		help.AddToRootCmd(c)
	}

	flags.AddPersistentFlags(c, ec)
	_ = viper.BindPFlags(c.PersistentFlags())

	c.SetOut(ec.Stdout)
	c.SetErr(ec.Stderr)

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

	return c
}

// initConfig ensures that precedence of configuration setting is correct
// precedence:
// flag -> environment -> configuration file value -> flag default
func initConfig(appName string, cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(appName)
	v.AddConfigPath(".")
	if dir, err := os.UserConfigDir(); err == nil {
		v.AddConfigPath(filepath.Join(dir, "solas"))
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
	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
