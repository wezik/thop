package platform

import (
	"os"
	"os/exec"

	"thop/pkg/platform"
)

func SystemExit() platform.ExitFn {
	return os.Exit
}

func SystemGetwd() platform.GetwdFn {
	return os.Getwd
}

func SystemOpenFile() platform.OpenFileFn {
	return os.OpenFile
}

func SystemExec() platform.ExecFn {
	return Exec
}

func Exec(cmd *exec.Cmd) (string, int, error) {
	res, err := cmd.Output()
	return string(res), cmd.ProcessState.ExitCode(), err
}
