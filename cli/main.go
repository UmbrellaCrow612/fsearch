package main

import (
	"github.com/UmbrellaCrow612/fsearch/cli/args"
	"github.com/UmbrellaCrow612/fsearch/cli/runner"
)

func main() {
	argsMap := args.Parse()
	runner.Run(argsMap)
}
