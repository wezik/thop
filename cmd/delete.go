package cmd

import (
	"thop/internal/service"
	"thop/internal/types/project"

	"github.com/spf13/cobra"
)

func deleteCmd(appService service.Service) *cobra.Command {
	return &cobra.Command{
		Use:     "delete [project]",
		Short:   "Delete a tmux session/project",
		Aliases: []string{"d"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) == 0 {
				projectName = ""
			} else {
				projectName = args[0]
			}

			return appService.DeleteProject(project.Name(projectName))
		},
	}
}
