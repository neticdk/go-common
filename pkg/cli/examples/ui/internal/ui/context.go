package ui

import "github.com/neticdk/go-common/pkg/cli/cmd"

type Context struct {
	EC *cmd.ExecutionContext
}

func NewContext() *Context {
	return &Context{}
}
