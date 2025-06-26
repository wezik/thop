package main

import (
	"fmt"
	"os"
	"thop/cmd"
	"thop/internal/config"
	"thop/internal/executor"
	"thop/internal/fsystem"
	"thop/internal/logger"
	"thop/internal/multiplexer"
	"thop/internal/selector"
	"thop/internal/service"
	"thop/internal/storage"
)

func main() {
	var logFile string

	rootCmd := cmd.GetRootCmd()
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Path to the log file")

	// a bit of manual DI
	c, err := config.New()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	if logFile == "" {
		logFile = c.GetConfigDir() + "/thop.log"
	}

	if err := logger.Init(logFile); err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	s := &storage.YamlStorage{Config: c, FileSystem: &fsystem.OsFileSystem{}}
	e := &executor.ShellExecutor{}
	cmd.AppService = &service.AppService{
		Selector: &selector.FzfProjectSelector{E: e},
		Multiplexer: &multiplexer.TmuxMultiplexer{
			ActiveTmuxSession: os.Getenv("TMUX"),
			Client:            &multiplexer.TmuxClientImpl{E: e},
		},
		Storage: s,
		Config:  c,
		E:       e,
	}

	cmd.Execute()
}
