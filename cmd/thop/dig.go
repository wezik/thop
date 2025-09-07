package main

import (
	"os"

	"thop/internal/domain/log"
	"thop/internal/domain/thop"

	logAdapters "thop/internal/adapters/log"
	multiplexerAdapters "thop/internal/adapters/multiplexer"
	platformAdapters "thop/internal/adapters/platform"
	selectorAdapters "thop/internal/adapters/selector"
	templateAdapters "thop/internal/adapters/template"

	"go.uber.org/dig"
)

// constructors provided to the container
func dependencies() []any {
	var groupedDependencies []any
	groupedDependencies = append(groupedDependencies, platformAdapters.SystemFunctions()...)
	dependencies := []any{
		thop.New,

		groupLogger,
		templateConfig,

		selectorAdapters.NewFzfSelector,
		templateAdapters.NewYamlStorage,
		multiplexerAdapters.NewTmuxMultiplexer,
	}

	return append(groupedDependencies, dependencies...)
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

func groupLogger() log.Logger {
	loggers := []log.Logger{}

	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	path := dir + "/thop/thop-new.log"

	fileLogger, err := logAdapters.NewFileLogger(logAdapters.LogFile(path), os.OpenFile)
	if err == nil {
		loggers = append(loggers, fileLogger)
	}

	loggers = append(loggers, logAdapters.NewEchoLogger(logAdapters.Debug))
	return logAdapters.NewGroupLogger(loggers)
}

func templateConfig() *templateAdapters.TemplateConfig {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return &templateAdapters.TemplateConfig{
		FileStoragePath: dir + "/thop/templates",
	}
}
