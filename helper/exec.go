package helper

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/gobackup/gobackup/logger"
	"github.com/google/uuid"
)

var (
	spaceRegexp = regexp.MustCompile(`[\s]+`)
)

// Exec cli commands
func Exec(command string, args ...string) (output string, err error) {
	return ExecWithStdio(command, false, args...)
}

func ExecWithStdio(command string, stdout bool, args ...string) (output string, err error) {
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	commandArgs := []string{}
	if len(commands) > 1 {
		commandArgs = commands[1:]
	}
	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = os.Environ()

	var stdErr bytes.Buffer
	var stdOut bytes.Buffer
	cmd.Stderr = &stdErr

	if stdout {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &stdOut
	}

	err = cmd.Run()
	if err != nil {
		logger.Debug(fullCommand, " ", strings.Join(commandArgs, " "))
		err = errors.New(stdErr.String())
	}
	output = strings.Trim(stdOut.String(), "\n")

	return
}

// Execute multiple line script with stdio
func ExecScriptWithStdio(script string, stdout bool) (string, error) {
	tmpFileName, _ := uuid.NewUUID()
	tmpFile := path.Join(os.TempDir(), tmpFileName.String())

	f, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	f.WriteString(script)
	err = os.Chmod(tmpFile, 0755)
	if err != nil {
		return "", err
	}

	defer f.Close()
	defer os.Remove(tmpFile)

	return ExecWithStdio("sh", stdout, tmpFile)
}

// Execute multiple line script
func ExecScript(script string) (output string, err error) {
	return ExecScriptWithStdio(script, false)
}
