package main

import (
	"thop/internal/log"
	"thop/internal/platform"
	"thop/internal/selector"
	"thop/internal/template"
	pkg_log "thop/pkg/log"
	pkg_platform "thop/pkg/platform"
	"thop/pkg/thop"

	"go.uber.org/dig"
)

// constructors provided to the container
func dependencies() []any {
	return []any{
		thop.New,

		groupLogger,

		selector.NewFzfSelector,
		template.NewFileSystemStorage,

		platform.SystemExit,
		platform.SystemGetwd,
		platform.SystemOpenFile,
	}
}

func autowireThop() *thop.Thop {
	container := dig.New()
	for _, dep := range dependencies() {
		container.Provide(dep)
	}

	var autowired *thop.Thop
	container.Invoke(func(t *thop.Thop) {
		autowired = t
	})

	return autowired
}

func groupLogger(openFile pkg_platform.OpenFileFn) pkg_log.Logger {
	loggers := []pkg_log.Logger{}

	fileLogger, err := log.NewFileLogger(log.LogFile("thop.log"), openFile)
	if err == nil {
		loggers = append(loggers, fileLogger)
	}

	loggers = append(loggers, log.NewEchoLogger(log.Debug))
	return log.NewGroupLogger(loggers)
}
