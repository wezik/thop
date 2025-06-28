package cmd

import (
	"github.com/spf13/cobra"
)

func rootCmd(openCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:           "thop",
		Short:         "Thop is a quick & lightweight tmux session/project manager",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return openCmd.RunE(cmd, args) // defaults to open command
		},
	}
}
