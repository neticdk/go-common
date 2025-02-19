package main

import (
	"os"

	"helloworld/cmd"
)

var version = "HEAD"

func main() {
	os.Exit(cmd.Execute(version))
}
