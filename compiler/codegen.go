package compiler

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"pengu-lang/templates"
)

func GenerateCode(service *Service, outputDir string, filename string) error {
	tmpl, err := template.New("service").Parse(templates.ServiceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, service); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	outPath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write generated file: %w", err)
	}

	return nil
}
