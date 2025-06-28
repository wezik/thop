package cmd

import (
	"thop/internal/service"
	"thop/internal/types/project"

	"github.com/spf13/cobra"
)

func openCmd(appService service.Service) *cobra.Command {
	return &cobra.Command{
		Use:     "open [project]",
		Short:   "Open a tmux session/project",
		Aliases: []string{"o", "select", "s"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) == 0 {
				projectName = ""
			} else {
				projectName = args[0]
			}

			return appService.OpenProject(project.Name(projectName))
		},
	}
}
