package main

import (
	"os"
	"thop/internal/log"
	"thop/internal/multiplexer"
	"thop/internal/platform"
	"thop/internal/selector"
	"thop/internal/template"

	pkg_log "thop/pkg/log"
	"thop/pkg/thop"

	"go.uber.org/dig"
)

// constructors provided to the container
func dependencies() []any {
	return []any{
		thop.New,

		groupLogger,
		templateConfig,

		selector.NewFzfSelector,
		template.NewYamlStorage,
		multiplexer.NewTmuxMultiplexer,

		platform.SystemExit,
		platform.SystemGetwd,
		platform.SystemOpenFile,
		platform.SystemExec,
		platform.SystemMkdirAll,
		platform.SystemReadDir,
		platform.SystemWriteFile,
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

func groupLogger() pkg_log.Logger {
	loggers := []pkg_log.Logger{}

	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	path := dir + "/thop/thop-new.log"

	fileLogger, err := log.NewFileLogger(log.LogFile(path), os.OpenFile)
	if err == nil {
		loggers = append(loggers, fileLogger)
	}

	loggers = append(loggers, log.NewEchoLogger(log.Debug))
	return log.NewGroupLogger(loggers)
}

func templateConfig() *template.TemplateConfig {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return &template.TemplateConfig{
		FileStoragePath: dir + "/thop/templates",
	}
}
