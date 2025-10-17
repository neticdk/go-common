package main

import (
	"helloworld/cmd"
	"os"
)

var version = "HEAD"

func main() {
	os.Exit(cmd.Execute(version))
}
