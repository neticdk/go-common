package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ecctx "github.com/neticdk/go-common/pkg/cli/context"
	"github.com/neticdk/go-common/pkg/cli/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	assert.Panics(t, func() {
		NewRootCommand(nil)
	})

	assert.Panics(t, func() {
		NewRootCommand(&ecctx.ExecutionContext{})
	})

	ec := &ecctx.ExecutionContext{
		AppName: "test",
	}
	cmd := NewRootCommand(ec).Build()
	assert.NotNil(t, cmd)
	assert.Equal(t, "test", cmd.Name())
}

func TestWithInitFunc(t *testing.T) {
	ec := newEC()

	initFuncCalled := false
	initFunc := func(cmd *cobra.Command, args []string) error {
		initFuncCalled = true
		return nil
	}

	cmd := NewRootCommand(ec).WithInitFunc(initFunc).Build()
	assert.NotNil(t, cmd)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.True(t, initFuncCalled)

	initFuncError := fmt.Errorf("init error")
	initFunc = func(cmd *cobra.Command, args []string) error {
		return initFuncError
	}
	cmd = NewRootCommand(ec).WithInitFunc(initFunc).Build()

	err = cmd.Execute()
	assert.ErrorIs(t, err, initFuncError)
}

func TestWithPersistentFlags(t *testing.T) {
	ec := newEC()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "", "test flag")

	cmd := NewRootCommand(ec).WithPersistentFlags(flags).Build()
	assert.NotNil(t, cmd)

	assert.NotNil(t, cmd.PersistentFlags().Lookup("test-flag"))
}

func TestBuild(t *testing.T) {
	ec := newEC()

	cmd := NewRootCommand(ec).Build()
	assert.NotNil(t, cmd)
}

func Test_initConfig(t *testing.T) {
	appName := "test"
	cmd := &cobra.Command{
		Use: appName,
	}

	flags.AddPersistentFlags(cmd, &ecctx.ExecutionContext{})

	err := initConfig(appName, cmd)
	assert.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	configDir := filepath.Join(homeDir, ".config", appName)
	err = os.MkdirAll(configDir, 0o755)
	assert.NoError(t, err)

	configFile := filepath.Join(configDir, fmt.Sprintf("%s.yaml", appName))
	err = os.WriteFile(configFile, []byte("verbose: true"), 0o644)
	assert.NoError(t, err)
	defer os.RemoveAll(configDir)

	err = initConfig(appName, cmd)
	assert.NoError(t, err)
}

func Test_bindFlags(t *testing.T) {
	v := viper.New()
	v.Set("test-flag", "test-value")

	cmd := &cobra.Command{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "", "test flag")
	cmd.Flags().AddFlagSet(flags)

	err := bindFlags(cmd, v)
	assert.NoError(t, err)

	assert.Equal(t, "test-value", cmd.Flags().Lookup("test-flag").Value.String())

	cmd = &cobra.Command{}
	flags = pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "default-value", "test flag")
	cmd.Flags().AddFlagSet(flags)
	_ = cmd.Flags().Set("test-flag", "flag-value")

	err = bindFlags(cmd, v)
	assert.NoError(t, err)

	assert.Equal(t, "flag-value", cmd.Flags().Lookup("test-flag").Value.String())
}

func TestGenDocsCommand(t *testing.T) {
	ec := newEC()

	cmd := GenDocsCommand(ec)

	assert.NotNil(t, cmd)
	assert.Equal(t, "gendocs", cmd.Name())
}

func newEC() *ecctx.ExecutionContext {
	return ecctx.NewExecutionContext("test", "test", "0.0.0", os.Stdin, os.Stdout, os.Stderr)
}
