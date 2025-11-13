package utils

import (
	"strings"

	"github.com/UmbrellaCrow612/fsearch/args"
)

func GetSizeMultipler(argsMap *args.ArgsMap) int64 {
	var sizeMultiplier int64 = 1
	switch strings.ToUpper(argsMap.SizeType) {
	case "KB":
		sizeMultiplier = 1024
	case "MB":
		sizeMultiplier = 1024 * 1024
	case "GB":
		sizeMultiplier = 1024 * 1024 * 1024
	case "B", "":
		sizeMultiplier = 1
	default:
		sizeMultiplier = 1
	}

	return sizeMultiplier
}
