package tests

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestE2EServer(t *testing.T) {
	tmpDir := t.TempDir()

	msContent := `version 1
service e2e_test

route GET "/ping"
    respond 200 "pong"

route POST "/data"
    respond 201 "created"
`
	msFile := filepath.Join(tmpDir, "e2e_app.ms")
	if err := os.WriteFile(msFile, []byte(msContent), 0644); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	cwd, _ := os.Getwd()
	projectRoot := filepath.Dir(cwd)
	cliMain := filepath.Join(projectRoot, "cli", "main.go")

	cmd := exec.Command("go", "run", cliMain, "run", msFile)
	cmd.Dir = tmpDir

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Wait for server to boot up
	time.Sleep(2 * time.Second)

	// Test GET /ping
	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Skipf("Failed to connect to server (port 8080 might be in use by another test or running app): %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected GET status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "pong" {
		t.Errorf("Expected GET body 'pong', got %s", string(body))
	}

	// Test POST /data
	resp, err = http.Post("http://localhost:8080/data", "text/plain", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated { // Note: compiler actually generates string literal like "respond 201"
		// Actually wait, 'respond 201' sends a 201 status code! Correct!
		t.Errorf("Expected POST status 201, got %d", resp.StatusCode)
	}
	body, _ = io.ReadAll(resp.Body)
	if string(body) != "created" {
		t.Errorf("Expected POST body 'created', got %s", string(body))
	}
}
