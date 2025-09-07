package main

import (
	"os"

	"thop/internal/domain/log"
	"thop/internal/domain/thop"
	"thop/internal/adapters/platform"

	logAdapters "thop/internal/adapters/log"
	multiplexerAdapters "thop/internal/adapters/multiplexer"
	selectorAdapters "thop/internal/adapters/selector"
	templateAdapters "thop/internal/adapters/template"

	"go.uber.org/dig"
)

// constructors provided to the container
func dependencies() []any {
	return []any{
		thop.New,

		workingDirectory,
		groupLogger,
		templateConfig,

		selectorAdapters.NewFzfSelector,
		templateAdapters.NewYamlStorage,
		multiplexerAdapters.NewTmuxMultiplexer,

		platform.NewSystemPlatform,
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

func workingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func groupLogger(
	platform platform.Platform,
) log.Logger {
	return logAdapters.NewGroupLogger([]log.Logger{
		fileLogger(platform),
		logAdapters.NewEchoLogger(logAdapters.Debug),
	})
}

func fileLogger(
	platform platform.Platform,
) log.Logger {
	path := "/tmp/thop-new.log"
	fileLogger, err := logAdapters.NewFileLogger(logAdapters.LogFile(path), platform)
	if err != nil {
		panic(err)
	}
	return fileLogger
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
