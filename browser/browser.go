package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Open(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Run()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Run()
	case "darwin":
		return exec.Command("open", url).Run()
	}

	return fmt.Errorf("unsupported platform")
}
