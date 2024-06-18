package git

import (
	"log"
	"os/exec"
	"strings"
)

func RootDirectory() string {
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}
