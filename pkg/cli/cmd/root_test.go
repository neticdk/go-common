package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
		NewRootCommand(&ExecutionContext{})
	})

	ec := &ExecutionContext{
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

	c := NewRootCommand(ec).WithInitFunc(initFunc).Build()
	assert.NotNil(t, c)

	err := c.Execute()
	assert.NoError(t, err)
	assert.True(t, initFuncCalled)

	initFuncError := fmt.Errorf("init error")
	initFunc = func(cmd *cobra.Command, args []string) error {
		return initFuncError
	}
	c = NewRootCommand(ec).WithInitFunc(initFunc).Build()

	err = c.Execute()
	assert.ErrorIs(t, err, initFuncError)
}

func TestWithPersistentFlags(t *testing.T) {
	ec := newEC()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "", "test flag")

	c := NewRootCommand(ec).WithPersistentFlags(flags).Build()
	assert.NotNil(t, c)

	assert.NotNil(t, c.PersistentFlags().Lookup("test-flag"))
}

func TestBuild(t *testing.T) {
	ec := newEC()

	c := NewRootCommand(ec).Build()
	assert.NotNil(t, c)
}

func Test_initConfig(t *testing.T) {
	appName := "test"
	c := &cobra.Command{
		Use: appName,
	}

	AddPersistentFlags(c, &ExecutionContext{})

	err := initConfig(appName, c)
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

	err = initConfig(appName, c)
	assert.NoError(t, err)
}

func Test_bindFlags(t *testing.T) {
	v := viper.New()
	v.Set("test-flag", "test-value")

	c := &cobra.Command{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "", "test flag")
	c.Flags().AddFlagSet(flags)

	err := bindFlags(c, v)
	assert.NoError(t, err)

	assert.Equal(t, "test-value", c.Flags().Lookup("test-flag").Value.String())

	c = &cobra.Command{}
	flags = pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "default-value", "test flag")
	c.Flags().AddFlagSet(flags)
	_ = c.Flags().Set("test-flag", "flag-value")

	err = bindFlags(c, v)
	assert.NoError(t, err)

	assert.Equal(t, "flag-value", c.Flags().Lookup("test-flag").Value.String())
}

func TestGenDocsCommand(t *testing.T) {
	ec := newEC()

	c := GenDocsCommand(ec)

	assert.NotNil(t, c)
	assert.Equal(t, "gendocs", c.Name())
}

func newEC() *ExecutionContext {
	return NewExecutionContext("test", "test", "0.0.0")
}
