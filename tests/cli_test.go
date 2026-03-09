package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLI(t *testing.T) {
	cwd, _ := os.Getwd()
	projectRoot := filepath.Dir(cwd)
	cliMain := filepath.Join(projectRoot, "cli", "main.go")

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "pengu.exe")

	buildCmd := exec.Command("go", "build", "-o", binPath, cliMain)
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	t.Run("init command", func(t *testing.T) {
		cmd := exec.Command(binPath, "init")
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("init failed: %v", err)
		}

		if _, err := os.Stat(filepath.Join(tmpDir, "examples", "auth.ms")); os.IsNotExist(err) {
			t.Errorf("Expected examples/auth.ms to be created")
		}
	})

	t.Run("generate command", func(t *testing.T) {
		msContent := `version 1
service gen_test
route GET "/hi"
    respond 200 "hi"
`
		msFile := filepath.Join(tmpDir, "test.ms")
		os.WriteFile(msFile, []byte(msContent), 0644)

		cmd := exec.Command(binPath, "generate", "test.ms")
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		if _, err := os.Stat(filepath.Join(tmpDir, "generated", "test.go")); os.IsNotExist(err) {
			t.Errorf("Expected generated/test.go to be created")
		}
	})
}
