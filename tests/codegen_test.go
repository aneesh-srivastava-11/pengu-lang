package tests

import (
	"os"
	"path/filepath"
	"pengu-lang/compiler"
	"strings"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	service := &compiler.Service{
		Version: "1",
		Name:    "test_codegen",
		Routes: []compiler.Route{
			{
				Method: "GET",
				Path:   "/ping",
				Actions: []compiler.Action{
					{Type: "respond", Args: []string{"200", "pong"}},
				},
			},
		},
	}

	tmpDir := t.TempDir()
	outName := "main.go"

	err := compiler.GenerateCode(service, tmpDir, outName)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	outPath := filepath.Join(tmpDir, outName)
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	code := string(content)
	if !strings.Contains(code, "package main") {
		t.Errorf("Generated code missing package main")
	}
	if !strings.Contains(code, "test_codegen") {
		t.Errorf("Generated code missing service name")
	}
}
