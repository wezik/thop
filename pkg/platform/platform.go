package platform

import (
	"os"
)

type ExitFn func(code int)
type GetwdFn func() (string, error)
type OpenFileFn func(name string, flag int, perm os.FileMode) (*os.File, error)
