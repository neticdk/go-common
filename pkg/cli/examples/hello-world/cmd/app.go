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
