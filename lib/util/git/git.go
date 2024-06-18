package git

import (
	"os/exec"
	"strings"
)

func RootDirectory() string {
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	output, _ := command.Output()
	return strings.TrimSpace(string(output))
}
