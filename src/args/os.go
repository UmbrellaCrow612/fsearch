package args

import "runtime"

func getValidViewersForOS() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"notepad", "wordpad", "code", "sublime_text", "notepad++", "vim", "nano", "less", "more",
		}
	case "darwin": // macOS
		return []string{
			"open", "code", "subl", "vim", "nano", "less", "more", "textedit",
		}
	case "linux":
		return []string{
			"xdg-open", "code", "gedit", "kate", "nano", "vim", "less", "more", "nvim",
		}
	default:
		return []string{"code", "vim", "nano", "less", "more"}
	}
}
