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

type Thop struct {
	service service.Service
	rootCmd *cobra.Command
}

func New(s service.Service) *Thop {
	thop := &Thop{service: s}

	openCmd := openCmd(s)
	rootCmd := rootCmd(openCmd)

	rootCmd.AddCommand(createCmd(s))
	rootCmd.AddCommand(deleteCmd(s))
	rootCmd.AddCommand(editCmd(s))
	rootCmd.AddCommand(killCmd(s))
	rootCmd.AddCommand(openCmd)

	thop.rootCmd = rootCmd
	return thop
}

func (t *Thop) GetLogFileFlag() string {
	var logFile string
	t.rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Path to the log file")
	return logFile
}

// returns exit code and error
// 0 - success
// 1 - error
// 2 - cancelled
func (t *Thop) Run() (int, error) {
	logger.Info(fmt.Sprintf("Starting Thop with args %s", os.Args[1:]))

	if err := t.rootCmd.Execute(); err != nil {
		switch err := err.(type) {

		case problem.Problem:
			// special case for selector cancellation
			if selector.ErrSelectorCancelled.Equal(err) {
				logger.Message("Selection cancelled")
				return 2, nil
			}
		}

		logger.Error(err)
		return 1, err
	}

	return 0, nil
}
