package main

import (
	"os"

	"hello-world/cmd"
)

var version = "HEAD"

func main() {
	os.Exit(cmd.Execute(version))
}
