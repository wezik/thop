package main

//go:generate go run ./mock-generator.go

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	mocks := []string{
		"../internal/domain/log/logger.go",
		"../internal/domain/selector/selector.go",
		"../internal/domain/multiplexer/multiplexer.go",
		"../internal/domain/template/storage.go",
		"../internal/adapters/platform/system.go",
	}

	for _, m := range mocks {
		mockFile := "../test/gen/mock/" + "mock_" + strings.ReplaceAll(strings.TrimPrefix(m, "../internal/"), "/", "_")

		cmd := exec.Command(
			"go", "run", "go.uber.org/mock/mockgen@latest",
			"-source", m,
			"-destination", mockFile,
			"-package", "mock",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		log.Printf("Generating mock: %s -> %s\n", m, mockFile)
		if err := cmd.Run(); err != nil {
			log.Fatalf("failed to generate mock for %s: %v", m, err)
		}
	}

	log.Println("All mocks generated successfully!")
}
