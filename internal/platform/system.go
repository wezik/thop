package platform

import (
	"os"

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
