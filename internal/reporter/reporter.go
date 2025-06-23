package reporter

import (
	"os"
	"os/exec"
	"time"
)

var URL = "https://github.com/alirezaarzehgar/iflashc/issues/new"

func appearsSuccessful(cmd *exec.Cmd, timeout time.Duration) bool {
	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		return true
	case err := <-errc:
		return err == nil
	}
}

func Open() bool {
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	if os.Getenv("DISPLAY") != "" {
		cmds = append(cmds, []string{"xdg-open"})
	}
	cmds = append(cmds,
		[]string{"firefox"},
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"chromium"},
	)

	for _, args := range cmds {
		cmd := exec.Command(args[0], append(args[1:], URL)...)
		if cmd.Start() == nil && appearsSuccessful(cmd, 3*time.Second) {
			return true
		}
	}
	return false
}
