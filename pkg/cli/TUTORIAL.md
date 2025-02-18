# CLI packages Tutorial

## Introduction

The `cli` packages helps you build CLIs that conform to [Srcoll
11](https://scr.k8s.netic.dk/0011/).

It consists of a collection of small packages:

- `cmd` provides helpers for building the root command and sub-commands.
   It sets sensible defaults, global flags, configuration files, etc. It
   also adds an interface to sub-commands that makes completing, validating
   and running commands more uniform.
- `errors` provides error handling and error types.
- `ui` provides UI elements such as tables, spinners, select boxes, prompts,
  etc.

This tutorial covers basic usage and the core concepts of these packages.

### Target Audience

The target audience for this tutorial is someone proficient with go but not
necessarily an expert. You should expect to read up on how
[spf13/cobra](https://github.com/spf13/cobra) works.

### Prerequisites

- go 1.22+

## Getting Started

Let's create a hello-world CLI that prints a greeting when we run it.

First, initialize the project:

```bash
mkdir hello-world
cd hello-world
go mod init hello-world
```

Download the go-common module:

```bash
go get -u github.com/neticdk/go-common
```

Create the `cmd` directory which will hold the code for the CLI commands:

```bash
mkdir cmd
```

Note, as this is a project with a single executable, we will put command code in
`cmd/` and keep a minimal `main.go` in the root project directory. For more
advanced tools that has multiple executables, you should create a directory for
each executable named `cmd/<executable-name>/` with its own `main.go` file. This
is an advanced topic and is not covered here.

Create `cmd/app.go`:

```go
package cmd

import "github.com/neticdk/go-common/pkg/cli/cmd"

const (
        AppName   = "hello-world"
        ShortDesc = "A greeting app"
        LongDesc  = `This application greets the user with a friendly messages`
)

type AppContext struct {
        EC *cmd.ExecutionContext
}

func NewAppContext() *AppContext {
        return &AppContext{}
}
```

Here we set up the `AppContext` (Application Context). Think of the
Application Context as the container for application information and
dependencies. It will vary from project to project.

For now it just holds a pointer to the `ExecutionContext` which you will learn
more about later. Also notice the constants `AppName`, `ShortDesc` and
`LongDesc`. These are used to identify and document your CLI, i.e when using
`--help` or the `help` command.

Next, create `cmd/root.go`:

```go
package cmd

import (
    "os"

    "github.com/neticdk/go-common/pkg/cli/cmd"
    "github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd(ac *AppContext) *cobra.Command {
    c := cmd.NewRootCommand(ac.EC).
        Build()

    c.AddCommand(
        HelloCmd(ac),
    )

    return c
}

// Execute runs the root command and returns the exit code
func Execute(version string) int {
    ec := cmd.NewExecutionContext(
        AppName,
        ShortDesc,
        version,
        os.Stdin,
        os.Stdout,
        os.Stderr)
    ac := NewAppContext()
    ac.EC = ec
    ec.LongDescription = LongDesc
    rootCmd := NewRootCmd(ac)
    err := rootCmd.Execute()
    _ = ec.Spinner.Stop()
    if err != nil {
        ec.ErrorHandler.HandleError(err)
        return 1
    }
    return 0
}

```

This is be the main entry point for handling CLI commands. The `Execute()`
function sets up the `ExecutionContext` and `AppContext`, runs the root
command (`NewRootCmd().Execute()`) and handles errors.

The `NewRootCmd()` function adds a sub-command, `HelloCmd()`, but we haven't
set that up yet, so let's do that now.

Create `cmd/hello.go`:

```go
package cmd

import (
    "context"

    "github.com/neticdk/go-common/pkg/cli/cmd"
    "github.com/neticdk/go-common/pkg/cli/ui"
    "github.com/spf13/cobra"
)

func HelloCmd(ac *AppContext) *cobra.Command {
    o := &helloOptions{}
    c := cmd.NewSubCommand("hello", o, ac).
        WithShortDesc("Say hello!").
        Build()

    return c
}

type helloOptions struct {
    who string
}

func (o *helloOptions) Complete(_ context.Context, ac *AppContext) {
    if len(ac.EC.CommandArgs) > 0 {
        o.who = ac.EC.CommandArgs[0]
    } else {
        o.who = "World"
    }
}

func (o *helloOptions) Validate(_ context.Context, _ *AppContext) error { return nil }

func (o *helloOptions) Run(_ context.Context, ac *AppContext) error {
    ui.Info.Printf("Hello, %s!\n", o.who)
    return nil
}

```

What you created here is a sub-command named `hello`. It has a short
description, "Say hello!" and a struct `helloOptions` to hold information such
as flag values or in this case who to say hello to.

You may notice the `AppContext` being passed to `NewSubCommand()` along with
`helloOptions`. There is some magic going on here which will be explained later,
but for now you need to know that it makes sure that `Complete()`, `Validate()`
and `Run()` are executed when the sub-command runs.

The `Complete()` command checks if there is an additional argument given to the
sub-command and uses that argument as the value for who to greet. It defaults to
"World" if none is given.

The `Run()` command executes the command and in this case prints out the
message.

Now all you need to do it make sure you CLI runs.

Create `main.go` which will run `cmd.Execute()`:

```go
package main

import (
        "os"

        "hello-world/cmd"
)

var version = "HEAD"

func main() {
        os.Exit(cmd.Execute(version))
}
```

Install dependencies:

```bash
go mod tidy
```

The directory structure should look like this:

```console
.
├── cmd
│   ├── app.go
│   ├── hello.go
│   └── root.go
├── go.mod
├── go.sum
└── main.go
```

Now run the CLI:

```bash
go run .
```

That prints out:

```console
This application greets the user with a friendly messages

## Usage

hello-world [command] [flags]
hello-world [command]

Basic Commands:

Other Commands:
completion  Generate the autocompletion script for the specified shell
help        Help about any command

## Additional Commands

  hello       Say hello!

## Flags

  -d, --debug               Debug mode
      --log-format string   Log format (plain|json) (default "plain")
      --log-level string    Log level (debug|info|warn|error) (default "info")
      --no-color            Do not print color
  -h, --help                help for hello-world
  -v, --version             version for hello-world

Use "hello-world [command] --help" for more information about a command.
```

That doesn't say hello anything. That's because you didn't specify the `hello`
command yet. Let's add that:

```bash
go run . hello
```

Your application should greet you with a friendly 'Hello, World!' message:

```console
 INFO  Hello, World!
```

Running it with an argument:

```bash
go run . hello John
```

… prints a greeting to John:

```console
 INFO  Hello, John!
```

And that concludes the getting started guide. Next up is core concepts.

## Core Concepts

### The Execution Context

The `ExecutionContext` is a struct containing information relevant to the
execution of the CLI but not necessarily coupled to the application directly.

Some of the important things it contains:

- cobra related attributes such as `Command`, `CommandArgs` and command
  information such as name and descriptions
- general attributes such as `Logger`, `ErrorHandler` and `OutputFormat`
- persistent/global flags (`PFlags`) configuration

You instantiate it using:

```go
import "github.com/neticdk/go-common/pkg/cli/cmd"

// cmd.NewExecutionContext(appName, shortDesc, version string, stdin io.Reader, stdout, stderr io.Writer)
ec := cmd.NewExecutionContext("driver", "app for driving cars", "0.0.1", os.Stdin, os.Stdout, os.Stderr)
```

If the I/O arguments are `nil` the `ExecutionContex` sets values to the
default `os` pipes (`stdin`, `stdout`, `stderr`).

#### Using the `ExecutionContext`

The `ExecutionContext` can be used by itself but works better in most cases if
embedded in other context structs such an application context:

```go
type AppContext struct {
    EC *cmd.ExecutionContext

    // Example
    DB db.DB
}

ac := &AppContext{
    EC: cmd.NewExecutionContext(...)
}
```

Create it and pass it to `NewRootCmd` from `Execute()` in `cmd/root.go`:

```go
func Execute(version string) {
// ..
ec := cmd.NewExecutionContext(...)
ac := &AppContext{EC: ec}
rootCmd := NewRootCmd(ac)
```

The `ExecutionContext` is passed on to the functions used to create commands:

```go
func NewRootCmd(ac *AppContext) {
    c := cmd.NewRootCommand(ac.EC)
    // ...
}

// or

func NewSubCmd(ac *AppContext) {
    o := &options{}
    c := cmd.NewSubCommand("name", o, ac)
    // ...
}
```

In the root command the `ExecutionContext` is available for all fields of the
`*cobra.Command` passed to or returned from `NewRootCommand`:

```go
func NewRootCmd(ac *AppContext) {
    c := cmd.NewRootCommand(ac.EC).
        WithInitFunc(func(_ *cobra.Command, _ []string) error {
            ac.SetupDefaultGithubClient()
            ac.SetupDefaultGitRepository()
            return nil
        }).
        Build()
    c.PreRun = func(_ *cobra.Command, args []string) {
        ec.Logger.Info("some message")
    }
    // ...
}
```

For sub-commands, the `NewSubCommand` passed down the context to the
`Complete`, `Validate` and `Run` functions:

```go
func InitComponentCmd(ac *AppContext) *cobra.Command {
    o := &options{}
    c := cmd.NewSubCommand("command", o, ac).
        Build()
    // ...
}

// passing ac to NewSubCommand automatically makes it available as the second
// argument to Complete, Validate and Run

func (o *options) Complete(ctx context.Context, ac *AppContext) {
    ac.EC.Logger.Info("some info")
}

func (o *options) Validate(ctx context.Context, ac *AppContext) error {
    return nil
}

func (o *options) Run(ctx context.Context, ac *AppContext) error {
    err := ac.GithubClient.GetRepository(...)
}
```

Note that the third argument passed to `NewSubCommand` is generic, so anything
you pass there will end up becoming the second argument to the three functions.
More on this later.

### Persistent/Global Flags

The `cmd` package comes with the default persistent flags, some of which are
permanent and some of which can be toggled. Persistent flags are always present
for all commands.

To enable a flag, set `<FLAG>Enabled` to `true`:

```go
ec := cmd.NewExecutionContext(...)
ec.PFlags.DryRunEnabled = true
```

See [`cmd/flags.go`](context/flags.go) for more information about available flags.

Flags, that can't be disabled:

- `--log-format`
- `--log-level`
- `--no-color`
- `--debug|-d`

#### Setting the output format

`ExecutionContext` stores the output format as an ENUM. To set it, enable the
supported format flags. For now, these are `--plain`, `--json`, `--yaml`,
and `--markdown` and they are mutually exclusive:

```go
ec.PFlags.PlainEnabled = true
ec.PFlags.JSONEnabled = true
ec.PFlags.YAMLEnabled = true
ec.PFlags.MarkdownEnabled = true
```

#### Enabling flags

A good place to set global flags is in the `Execute()` command in
`cmd/root.go`:

```go
// Execute runs the root command and returns the exit code
func Execute(version string) int {
    ec := cmd.NewExecutionContext(
        AppName,
        ShortDesc,
        version,
        os.Stdin,
        os.Stdout,
        os.Stderr)
    ac := NewAppContext()
    ec.PFlags.ForceEnabled = true
    ec.PFlags.JSONEnabled = true
    ac.EC = ec
    ec.LongDescription = LongDesc
    rootCmd := NewRootCmd(ac)
    err := rootCmd.Execute()
    _ = ec.Spinner.Stop()
    if err != nil {
        ec.ErrorHandler.HandleError(err)
        return 1
    }
    return 0
}
```

### Logging

The package uses `log/slog` for logging.

#### Default Handler

It includes a default handler using the `pterm` package. The default format
depends on whether there is a TTY attached. If there is then the format is
`plain` text unless `--log-format` is used. Otherwise is uses `json`.

The default log level is INFO. To change this, use the `--log-level` flag.

#### Using your own handler

You can change the handler after initialing the context:

```go
ec := cmd.NewExecutionContext(...)

handler := ...
ec.Logger = slow.New(handler)

```

#### Using the logger

Just use `ec.Logger.Info()` (or `ac.EC.Logger.Info()`) etc. like you would use `slog`.

### Error Handling

The `ExecutionContext` comes with a default error handler
(`errors.DefaultHandler`) that implements this interface:

```go
type Handler interface {
    HandleError(err error)
    NewGeneralError(message, helpMsg string, err error, code int) *GeneralError
    SetWrap(wrap bool)
    SetWrapWidth(width int)
}
```

Errors bubble up from your application and are handled in the `Execute()`
function in `cmd/root.go`:

```go
func Execute(...) int {
    // ...
    err := rootCmd.Execute()
    _ = ec.Spinner.Stop()
    if err != nil {
        ec.ErrorHandler.HandleError(err)
        return 1
    }
    return 0
}
```

The `DefaultHandler` handles two types of errors:

- `ErrorWithHelp` is used for printing user friendly errors with context
- The built in `error` handling all other errors

`ErrorWithHelp` is an interface:

```go
type ErrorWithHelp interface {
    error
    Help() string
    Unwrap() error // Optional: for wrapped errors
    Code() int     // Optional: for error codes
}
```

So using any type that implements this interface  as an `error` will make sure
it is printed out in a nice format for the user.

There are two error types included that implements this interface:

- `GeneralError` - which can handle most cases
- `InvalidArgumentError` - specifically made for parsing arguments and flags

Use `GeneralError` like this:

```go
import "github.com/neticdk/go-common/pkg/cli/errors"

const ErrorCodeParsingError = 42

func myFunc() error {
    // ...
    return &errors.GeneralError{
        Message: "Could not parse config.json",
        HelpMsg: "This happens when the file format is invalid. See details for more.",
        Err:     err,
        CodeVal: ErrorCodeParsingError,
    }
}
```

There is also a short hand helper for `GeneralError`:

```go
return ac.EC.ErrorHandler.NewGeneralError(
    "Could not parse config.json",
    "This happens...",
    err,
    0)
```

Using error codes is optional.

Use `InvalidArgumentError` like this:

```go
import "github.com/neticdk/go-common/pkg/cli/errors"

func (o *options) Validate(ctx context.Context, ac *AppContext) {
    return &errors.InvalidArgumentError{
        Flag:    "name",
        Val:     o.Name,
        Context: "It must be an ASCII string of minimum 3 characters length.",
    }
    // ...
}
```

### Commands

The two helpers `cmd.NewRootCommand` and `cmd.NewSubCommand` are used to create
root commands and sub-commands respectively.

They implement the builder pattern (chaining functions that modifies the return
value) and both return a `*cobra.Command`. This means that everything you can do
with Cobra, you can do with these helpers. They are just meant to set reasonable
defaults and enable some features that you almost always want.

`cmd.NewRootCommand(ec *cmd.ExecutionContext).Build()` creates a root command with:

- usage, descriptions, versions, etc. set
- default global flags added
- logging enabled
- configuration through configuration files, environment variables, and
  flags enabled (via `viper`)
- colors added to `help` commands
- two command groups added (`cmd.GroupBase`, and `cmd.GroupOther`)
- a hidden command `gendocs` for generating documentation

`cmd.NewSubCommand(name string, runner cmd.SubCommandRunner[*AppContext],
runnerArg *AppContext)).Build()` works a little different. It takes the name of
the command, a runner interface and an argument passed to the functions of the
runner interface as arguments. Let's break the last two down because they are
important to understand.

The runner interface looks like this:

```go
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
```

This means that the type must implement these three functions. This pattern is
called the complete-validate-run pattern and is used by the kubernetes project
amongst others. It adds predictability at the cost of some flexibility. Behind
the scenes it overrides the `RunE` field of `*cobra.Command` with a function
that looks like this:

```go
return func(cmd *cobra.Command, args []string) error {
    ctx := cmd.Context()
    runner.Complete(ctx, runnerArg)
    if err := runner.Validate(ctx, runnerArg); err != nil {
        return err
    }
    return runner.Run(ctx, runnerArg)
}
```

You can always set `RunE` on `cobra.Cobra` to do your own thing.

When the sub-command runs, it runs these three functions in order:

Use `Complete` to complete any settings/configuration/flags/etc. before
validation. It doesn't return anything.

Given this struct:

```go
type options struct {
    name string
    age  int
    car  string
    dest string
}
```

The `Complete` function could do something like this:

```go
func (o *options) Complete(ctx context.Context, ac *AppContext) {
    if o.age > 50 {
        o.car == "Mercedes"
    }
}
```

Use `Validate` to validate flags, arguments and other things. It always
returns error.

```go
func (o *options) Validate(ctx context.Context, ac *AppContext) error {
    if o.age < 18 {
        return ac.EC.ErrorHandler.NewGeneralError(
            "Child detected",
            "The person under under 18 years old and cannot drive a car."
            err,
            0)
    }
}
```

Finally, use the `Run` function to run the command. It also returns error:

```go
func (o *options) Run(ctx context.Context, ac *AppContext) error {
    return car.Drive(ctx, o.name, o.car, o.dest)
}
```

Notice that `ac *AppContext` in the function signatures? That is the third
argument passed to `NewSubCommand`. This makes that argument available to the
rest of the application.

To fulfill the `SubCommandRunner` interface the following must be in place:

- All three functions must be implemented on the struct passed as the
  second argument to `NewSubCommand`.
- The type of the second argument to each of the three functions must be the
  same as the type of the third argument passed to `NewSubCommand`.

At this point you might be wondering how the options are populated from command
line flags in the first place? This is where binding comes in:

```go
func DriveCmd(ac *AppContext) *cobra.Command {
    o := &driveOptions{}
    c := cmd.NewSubCommand("drive", o, ac).
        WithShortDesc("Drive a car").
        WithLongDesc(driveCmdLongDescription()).
        WithExample(driveCmdExample()).
        WithGroupID(groupComponent).
        Build()

    o.bindFlags(c.Flags())
    c.Flags().SortFlags = false
    _ = c.MarkFlagRequired("name")
    _ = c.MarkFlagRequired("age")
    return c
}
```

Notice the `bindFlags` function. It binds the command line flags to the options.
Let's see an example of how to implement it:

```go
import "github.com/spf13/pflag"

func (o *driveOptions) bindFlags(f *pflag.FlagSet) {
    f.StringVar(&o.name, "name", "", "Driver name")
    f.IntVar(&o.age, "age", 0, "Driver age")
}
```

This binds the struct fields `o.name` and `o.age` to the values passed to
`--name` and `--age` respectively.

Now you may also have noticed that the example above used some builder
functions. They are essentially wrappers to set fields on the `cobra.Command`
struct. You may use them or just set the field values yourself after creating
the sub-command.

#### Accessing the command

The `cobra.Command` can be accessed through `ExecutionContext.Command`.

#### Accessing the command args

The `args` can be accessed through `ExecutionContext.CommandArgs`.

### UI Elements

There's a couple of UI elements included in the `ui` package. There are:

- `ui.NewTable()` - for creating tables
- `ui.Select()` - for creating selection inputs
- `ui.Spin()` - for running functions with a spinner
- `ui.Confirm()` - for creating confirmation prompts
- `ui.Prompt()` - for prompting for input
- a range of prefix writers (`ui.Info`, `ui.Success`, etc.) each of
  which are chained command that takes printer functions (e.g.
  `ui.Success.Println("Yay!")`

Look at the [package](ui/) to see what is available.

## Building a Simple Project

TBD
