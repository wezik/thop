package cmd

import (
	"thop/internal/service"
	"thop/internal/types/project"

	"github.com/spf13/cobra"
)

func killCmd(appService service.Service) *cobra.Command {
	return &cobra.Command{
		Use:     "kill [session]",
		Short:   "Kill active tmux session",
		Aliases: []string{"k"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) == 0 {
				projectName = ""
			} else {
				projectName = args[0]
			}

			return appService.KillSession(project.Name(projectName))
		},
	}
}
