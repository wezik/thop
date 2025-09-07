package template

import (
	"os"
	"strings"
	"thop/pkg/template"
)

func mapToTemplate(yamlTemplate *YamlTemplate, filePath template.FilePath) *template.Template {
	var windows []*template.Window
	for _, window := range yamlTemplate.Windows {
		windows = append(windows, mapToWindow(window))
	}
	return template.NewTemplate(
		filePath,
		template.Version(yamlTemplate.Version),
		yamlTemplate.Name,
		yamlTemplate.SessionName,
		yamlTemplate.Path,
		yamlTemplate.Commands,
		windows,
	)
}

func mapToWindow(yamlWindow *YamlWindow) *template.Window {
	var panes []*template.Pane
	for _, pane := range yamlWindow.Panes {
		panes = append(panes, mapToPane(pane))
	}
	return template.NewWindow(
		yamlWindow.Name,
		yamlWindow.Path,
		yamlWindow.Layout,
		yamlWindow.Commands,
		panes,
	)
}

func mapToPane(yamlPane *YamlPane) *template.Pane {
	return template.NewPane(
		yamlPane.Active,
		yamlPane.Commands,
		yamlPane.Path,
	)
}

func mapToYamlTemplate(template *template.Template) *YamlTemplate {
	var windows []*YamlWindow
	for _, window := range template.Windows() {
		windows = append(windows, mapToYamlWindow(window))
	}
	return &YamlTemplate{
		Version:     int(template.Version()),
		Name:        template.Name(),
		SessionName: template.SessionName(),
		Path:        template.Path(),
		Commands:    template.Commands(),
		Windows:     windows,
	}
}

func mapToYamlWindow(window *template.Window) *YamlWindow {
	var panes []*YamlPane
	for _, pane := range window.Panes() {
		panes = append(panes, mapToYamlPane(pane))
	}
	return &YamlWindow{
		Name:     window.Name(),
		Path:     window.Path(),
		Layout:   window.Layout(),
		Commands: window.Commands(),
		Panes:    panes,
	}
}

func mapToYamlPane(pane *template.Pane) *YamlPane {
	return &YamlPane{
		Active:   pane.Active(),
		Commands: pane.Commands(),
		Path:     pane.Path(),
	}
}

func mapToTemplateFile(file os.DirEntry, root string) *template.File {
	return template.NewFile(file.Name(), root)
}

func fileSafeName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.TrimPrefix(name, "-")
	name = name + ".yaml"
	return name
}
