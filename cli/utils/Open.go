package utils

import (
	"fmt"
	"os/exec"

	"github.com/UmbrellaCrow612/fsearch/cli/args"
	"github.com/UmbrellaCrow612/fsearch/cli/shared"
)

func OpenMatchEntry(entry shared.MatchEntry) error {
	if entry.Entry.IsDir() {
		explorer := args.GetFileExplorerForOS()
		cmd := exec.Command(explorer, entry.Path)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to open folder %s: %v", entry.Path, err)
		}
		fmt.Printf("Opened folder %s with %s\n", entry.Path, explorer)
		return nil
	} else {
		viewers := args.GetValidViewersForOS()
		for _, viewer := range viewers {
			cmd := exec.Command(viewer, entry.Path)
			if err := cmd.Start(); err == nil {
				fmt.Printf("Opened file %s with %s\n", entry.Path, viewer)
				return nil
			}
		}
		return fmt.Errorf("could not open file %s with any viewer", entry.Path)
	}
}
