package main

import (
	"os"
	thop "thop/cmd"
	"thop/internal/config"
	"thop/internal/executor"
	"thop/internal/fsystem"
	"thop/internal/logger"
	"thop/internal/multiplexer"
	"thop/internal/selector"
	"thop/internal/service"
	"thop/internal/storage"

	"go.uber.org/dig"
)

func main() {
	container := buildContainer()

	var app *thop.Thop
	container.Invoke(func(t *thop.Thop) { app = t })

	logFile := app.GetLogFileFlag()
	initLogger(logFile, container)

	exitCode, _ := app.Run()
	os.Exit(exitCode)
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.NewConfig)
	container.Provide(fsystem.NewOsFileSystem)
	container.Provide(executor.NewShellExecutor)

	container.Provide(multiplexer.NewTmuxClientImpl)
	container.Provide(multiplexer.NewTmuxMultiplexer)

	container.Provide(selector.NewFzfProjectSelector)

	container.Provide(storage.NewYamlStorage)

	container.Provide(service.NewAppService)

	container.Provide(thop.New)

	return container
}

func initLogger(logFile string, container *dig.Container) error {
	var c *config.Config
	if err := container.Invoke(func(config *config.Config) {
		c = config
	}); err != nil {
		return err
	}

	if logFile == "" {
		logFile = c.GetConfigDir() + "/thop.log"
	}

	return logger.Init(logFile)
}
