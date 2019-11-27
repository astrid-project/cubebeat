package main

import (
	"os"

	"gitlab.com/astrid-repositories/wp2/cubebeat/cmd"

	_ "gitlab.com/astrid-repositories/wp2/cubebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
