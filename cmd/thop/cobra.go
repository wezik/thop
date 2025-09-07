package main

import (
	"thop/internal/domain/template"
	"thop/internal/domain/thop"

	"github.com/spf13/cobra"
)

func root(thop *thop.Thop) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "thop",
		Short: "A multiplexer session/project manager",
	}

	for _, sub := range subCommands(thop) {
		rootCmd.AddCommand(sub)
	}

	return rootCmd
}

func subCommands(thop *thop.Thop) []*cobra.Command {
	return []*cobra.Command{
		createTemplate(thop),
		deleteTemplate(thop),
		editTemplate(thop),
		killSession(thop),
		openTemplate(thop),
	}
}

func killSession(t *thop.Thop) *cobra.Command {
	var name string

	command := &cobra.Command{
		Use:     "kill",
		Short:   "Kill active session",
		Aliases: []string{"k", "kill"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if name != "" {
				search := thop.Search{
					Phrase: name,
				}
				return t.KillSearch(search)
			}
			return t.KillSelect()
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "Name of the template")

	return command
}

func createTemplate(t *thop.Thop) *cobra.Command {
	var name string
	var path string

	command := &cobra.Command{
		Use:     "create",
		Short:   "Create session template",
		Aliases: []string{"c", "a", "add", "new"},
		RunE: func(cmd *cobra.Command, args []string) error {
			createTemplate := template.CreateTemplate{
				Name: name,
				Path: path,
			}
			_, err := t.Create(createTemplate)
			return err
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "Name of the template")
	command.Flags().StringVarP(&path, "path", "p", "", "Path to template session root")

	return command
}

func deleteTemplate(t *thop.Thop) *cobra.Command {
	var name string

	command := &cobra.Command{
		Use:     "delete",
		Short:   "Delete session template",
		Aliases: []string{"d", "delete"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if name != "" {
				search := thop.Search{
					Phrase: name,
				}
				return t.DeleteSearch(search)
			}
			return t.DeleteSelect()
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "Name of the template")

	return command
}

func editTemplate(t *thop.Thop) *cobra.Command {
	var name string

	command := &cobra.Command{
		Use:     "edit",
		Short:   "Edit session template",
		Aliases: []string{"e", "edit"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if name != "" {
				search := thop.Search{
					Phrase: name,
				}
				return t.EditSearch(search)
			}
			return t.EditSelect()
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "Name of the template")

	return command
}

func openTemplate(t *thop.Thop) *cobra.Command {
	var name string

	command := &cobra.Command{
		Use:     "open",
		Short:   "Open session template or attach to active session",
		Aliases: []string{"o", "open"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if name != "" {
				search := thop.Search{
					Phrase: name,
				}
				return t.OpenSearch(search)
			}
			return t.OpenSelect()
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "Name of the template")

	return command
}
