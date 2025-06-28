package cmd

import (
	"thop/internal/service"
	"thop/internal/types/project"

	"github.com/spf13/cobra"
)

func editCmd(appService service.Service) *cobra.Command {
	return &cobra.Command{
		Use:     "edit [project]",
		Short:   "Edit a tmux session/project using $EDITOR",
		Aliases: []string{"e"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) == 0 {
				projectName = ""
			} else {
				projectName = args[0]
			}

			return appService.EditProject(project.Name(projectName))
		},
	}
}
