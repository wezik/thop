package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func listImports(filePath string) []string {
	f, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil
	}

	var result []string
	for _, imp := range f.Imports {
		result = append(result, imp.Path.Value)
	}
	return result
}

func listGoFilesFromTree(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		files = append(files, path)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func TestAdapterDependencies(t *testing.T) {
	// given
	files := listGoFilesFromTree("../../internal/adapters/")
	restrictedImportPaths := []string{
		"/cmd/",
		"\"github.com/spf13/cobra\"",
		"\"go.uber.org/dig\"",
	}

	type ImportEntry struct {
		ImportPath string
		FilePath   string
	}

	// when
	var dependencies []ImportEntry
	for _, file := range files {
		imports := listImports(file)
		for _, imp := range imports {
			dependencies = append(dependencies, ImportEntry{
				ImportPath: imp,
				FilePath:   file,
			})
		}
	}

	// then
	t.Run("no restricted dependencies", func(t *testing.T) {
		for _, dep := range dependencies {
			if slices.Contains(restrictedImportPaths, dep.ImportPath) {
				t.Errorf("%s depends on %s", dep.FilePath, dep.ImportPath)
			}
		}
	})

	t.Run("no external logger dependencies", func(t *testing.T) {
		allowedPath := "../../internal/adapters/log/"
		for _, dep := range dependencies {
			// skip logger implementations
			if strings.Contains(dep.FilePath, allowedPath) {
				continue
			}

			if strings.Contains(dep.ImportPath, "log") && dep.ImportPath != "\"thop/internal/domain/log\"" {
				t.Errorf("%s depends on external logger implementation %s", dep.FilePath, dep.ImportPath)
			}
		}
	})
}

func TestPkgDependencies(t *testing.T) {
	// given
	files := listGoFilesFromTree("../../internal/domain/")
	restrictedImportPaths := []string{
		"/cmd/",
		"/internal/adapters/",
		"\"github.com/spf13/cobra\"",
		"\"go.uber.org/dig\"",
	}
	type ImportEntry struct {
		ImportPath string
		FilePath   string
	}

	// when
	var dependencies []ImportEntry
	for _, file := range files {
		imports := listImports(file)
		for _, imp := range imports {
			dependencies = append(dependencies, ImportEntry{
				ImportPath: imp,
				FilePath:   file,
			})
		}
	}

	// then
	t.Run("no restricted dependencies", func(t *testing.T) {
		for _, dep := range dependencies {
			if slices.Contains(restrictedImportPaths, dep.ImportPath) {
				t.Errorf("%s depends on %s", dep.FilePath, dep.ImportPath)
			}
		}
	})

	t.Run("no external logger dependencies", func(t *testing.T) {
		for _, dep := range dependencies {
			if strings.Contains(dep.ImportPath, "log") && dep.ImportPath != "\"thop/internal/domain/log\"" {
				t.Errorf("%s depends on external logger implementation %s", dep.FilePath, dep.ImportPath)
			}
		}
	})
}


