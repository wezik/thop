package main

import (
	"fmt"
	"os"
	"thop/cmd"
	"thop/internal/config"
	"thop/internal/executor"
	"thop/internal/fsystem"
	"thop/internal/logger"
	"thop/internal/messenger"
	"thop/internal/multiplexer"
	"thop/internal/selector"
	"thop/internal/service"
	"thop/internal/storage"

	"github.com/spf13/cobra"
)

func main() {
	var logFile string

	rootCmd := cmd.GetRootCmd()
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Path to the log file")

	cobra.OnInitialize(func() {
		// a bit of manual DI
		c, err := config.New()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			os.Exit(1)
		}

		if logFile == "" {
			logFile = c.GetConfigDir() + "/thop.log"
		}

		logger, err := logger.New(logFile)
		if err != nil {
			fmt.Println("Failed to initialize logger:", err)
			os.Exit(1)
		}

		msg := &messenger.Messenger{Logger: logger}
		s := &storage.YamlStorage{Config: c, FileSystem: &fsystem.OsFileSystem{}}
		e := &executor.ShellExecutor{Logger: logger}
		cmd.AppService = &service.AppService{
			Selector: &selector.FzfProjectSelector{E: e},
			Multiplexer: &multiplexer.TmuxMultiplexer{
				ActiveTmuxSession: os.Getenv("TMUX"),
				Client:            &multiplexer.TmuxClientImpl{E: e},
				Messenger:         msg,
			},
			Storage: s,
			Config:  c,
			E:       e,
		}
	})

	cmd.Execute()
}
