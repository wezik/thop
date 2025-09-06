package architecture

import (
    "go/parser"
    "go/token"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestInternalDoesNotImportCmd(t *testing.T) {
    root := "../../internal/" // internal root
    fset := token.NewFileSet()

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() || !strings.HasSuffix(path, ".go") {
            return nil
        }

        f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
        if err != nil {
            return err
        }

        for _, imp := range f.Imports {
            if strings.Contains(imp.Path.Value, "/cmd/") {
                t.Errorf("%s imports cmd package %s", path, imp.Path.Value)
            }
        }

        return nil
    })
    if err != nil {
        t.Fatal(err)
    }
}

func TestPkgDoesNotImportInternal(t *testing.T) {
    root := "../../pkg/" // pkg root
    fset := token.NewFileSet()

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() || !strings.HasSuffix(path, ".go") {
            return nil
        }

        f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
        if err != nil {
            return err
        }

        for _, imp := range f.Imports {
            if strings.Contains(imp.Path.Value, "/internal/") {
                t.Errorf("%s imports internal package %s", path, imp.Path.Value)
            }
        }

        return nil
    })
    if err != nil {
        t.Fatal(err)
    }
}

func TestPkgDoesNotImportCmd(t *testing.T) {
    root := "../../pkg/" // pkg root
    fset := token.NewFileSet()

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() || !strings.HasSuffix(path, ".go") {
            return nil
        }

        f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
        if err != nil {
            return err
        }

        for _, imp := range f.Imports {
            if strings.Contains(imp.Path.Value, "/cmd/") {
                t.Errorf("%s imports cmd package %s", path, imp.Path.Value)
            }
        }

        return nil
    })
    if err != nil {
        t.Fatal(err)
    }
}
