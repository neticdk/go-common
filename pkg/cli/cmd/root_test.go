package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	ec := NewExecutionContext("test", "test", "0.0.0")
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

func TestWithUpdateChecker(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tag_name": "v1.2.0"}`))
	}))
	defer ts.Close()

	ec := newEC()
	ec.Version = "v1.0.0"
	var buf bytes.Buffer
	ec.Stderr = &buf
	ec.initLogger() // Re-bind the logger so it writes to the test buffer

	checker := NewUpdateChecker(ec, "owner", "repo", "")
	checker.githubURL = ts.URL
	checker.cacheDir = t.TempDir()
	checker.cacheDuration = 0

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tag_name": "v2.0.0"}`))
	}))
	defer ts2.Close()

	checker2 := NewUpdateChecker(ec, "owner2", "repo2", "",
		WithAppName("dep2"),
		WithCurrentVersion("v1.5.0"),
	)
	checker2.githubURL = ts2.URL
	checker2.cacheDir = t.TempDir()
	checker2.cacheDuration = 0

	c := NewRootCommand(ec).
		WithUpdateChecker(checker).
		WithUpdateChecker(checker2).
		Build()
	c.RunE = func(cmd *cobra.Command, args []string) error {
		time.Sleep(100 * time.Millisecond) // wait for async routine
		return nil
	}

	err := c.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "🚀 A new version is available! (v1.0.0 -> v1.2.0)")
	assert.Contains(t, buf.String(), "🚀 A new version is available! (v1.5.0 -> v2.0.0)")
}

func newEC() *ExecutionContext {
	return NewExecutionContext("test", "test", "0.0.0")
}
