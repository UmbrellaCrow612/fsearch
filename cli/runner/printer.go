package runner

import (
	"bufio"
	"os"
	"strings"

	"github.com/UmbrellaCrow612/fsearch/cli/args"
	"github.com/UmbrellaCrow612/fsearch/cli/out"
	"github.com/UmbrellaCrow612/fsearch/cli/shared"
)

func printMatchs(m []shared.MatchEntry, argMap *args.ArgsMap) {
	for _, i := range m {
		nameAndPath := i.Entry.Name() + " " + i.Path

		// Only print file content if type is "file" and Lines > 0
		if argMap.Lines > 0 && strings.ToLower(argMap.Type) == "file" && !i.Entry.IsDir() {
			contentPreview, err := readFirstNLines(i.Path, argMap.Lines)
			if err != nil {
				out.WriteToStderr(nameAndPath + " ERROR reading file: " + err.Error())
				continue
			}
			// Print name + path, then content separated by a newline
			out.WriteToStdout(nameAndPath + "\n" + contentPreview)
		} else {
			out.WriteToStdout(nameAndPath)
		}
	}
}

// Reads the first n lines of a file
func readFirstNLines(filePath string, n int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	count := 0
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		count++
		if count >= n {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}
