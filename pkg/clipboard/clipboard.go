package clipboard

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func CopyToClipboard(text string) error {
	var copyCmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		copyCmd = exec.Command("xclip", "-selection", "clipboard")
	case "darwin":
		copyCmd = exec.Command("pbcopy")
	case "windows":
		copyCmd = exec.Command("clip")
	default:
		return errors.New("unsupported operating system")
	}

	copyCmd.Stdin = strings.NewReader(text)
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return nil
}
