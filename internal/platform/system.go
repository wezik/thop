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
	return execImpl
}

func execImpl(cmd *exec.Cmd) (string, int, error) {
	res, err := cmd.Output()
	return string(res), cmd.ProcessState.ExitCode(), err
}

func SystemMkdirAll() platform.MkdirAllFn {
	return mkdirAllImpl
}

func mkdirAllImpl(path string) error {
	return os.MkdirAll(path, 0755)
}

func SystemReadDir() platform.ReadDirFn {
	return os.ReadDir
}
