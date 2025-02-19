package main

import (
	"os"

	"pokemon/cmd"
)

var version = "HEAD"

func main() {
	os.Exit(cmd.Execute(version))
}
