package main

import (
	"github.com/UmbrellaCrow612/fsearch/args"
	"github.com/UmbrellaCrow612/fsearch/runner"
)

func main() {
	argsMap := args.Parse()
	runner.Run(argsMap)
}
