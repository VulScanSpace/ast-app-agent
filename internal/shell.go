package internal

import (
	"bytes"
	"os/exec"
)

func RunShellCmd(name string, args ...string) (err error, output, errMsg string) {
	cmd := exec.Command(name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	errMsg = stderr.String()
	output = stdout.String()
	return
}
