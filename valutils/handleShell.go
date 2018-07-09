package valutils

import (
	"bytes"
	"os/exec"
)

// AppVersionCheck tries to execute a shell command with commandline argument
// "--version" and returns the resulting commandline output or the error if
// one occurs.
func AppVersionCheck(binpath string) (string, error) {
	cmd := exec.Command(binpath, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
