package main

import (
	"github.com/UmbrellaCrow612/fsearch/src/args"
	"github.com/UmbrellaCrow612/fsearch/src/runner"
)

func main() {
	argsMap := args.Parse()
	runner.Run(argsMap)
}
