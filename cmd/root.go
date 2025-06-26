package cmd

import (
	"fmt"
	"os"
	"thop/internal/logger"
	"thop/internal/problem"
	"thop/internal/selector"
	"thop/internal/service"

	"github.com/spf13/cobra"
)

var AppService service.Service
var rootCmd = &cobra.Command{
	Use:           "thop",
	Short:         "Thop is a quick & lightweight tmux session/project manager",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// No args, defaults to open command
			return openCmd.RunE(cmd, args)
		}

		return nil
	},
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func Execute() {
	logger.Info(fmt.Sprintf("Executing args: %s", os.Args[1:]))

	if err := rootCmd.Execute(); err != nil {
		switch err := err.(type) {

		case problem.Problem:
			// special case for selector cancellation
			if selector.ErrSelectorCancelled.Equal(err) {
				logger.Message("Selection cancelled")
				os.Exit(0)
			}
		}

		logger.Error(err)
		os.Exit(1)
	}
}
