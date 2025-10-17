package main

import (
	"os"
	"ui/cmd"
)

var version = "HEAD"

func main() {
	os.Exit(cmd.Execute(version))
}
