package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"pengu-lang/compiler"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		initProject()
	case "generate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: pengu.exe generate <file.ms>")
			os.Exit(1)
		}
		generateCode(os.Args[2])
	case "build":
		if len(os.Args) < 3 {
			fmt.Println("Usage: pengu.exe build <file.ms>")
			os.Exit(1)
		}
		buildBinary(os.Args[2])
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: pengu.exe run <file.ms>")
			os.Exit(1)
		}
		runService(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("pengu - Microservice DSL Compiler")
	fmt.Println("Commands:")
	fmt.Println("  init                  Create a starter project")
	fmt.Println("  generate <file.ms>    Generate Go code without building")
	fmt.Println("  build <file.ms>       Generate Go code and build binary")
	fmt.Println("  run <file.ms>         Generate Go code and run service")
}

func initProject() {
	os.MkdirAll("examples", 0755)

	authMs := `version 1
service auth

route POST "/login"
    log "login attempt"
    respond 200 "success"
`
	userMs := `version 1
service user

route GET "/profile"
    log "profile request"
    respond 200 "ok"
`
	os.WriteFile("examples/auth.ms", []byte(authMs), 0644)
	os.WriteFile("examples/user.ms", []byte(userMs), 0644)
	fmt.Println("Initialized new pengu project in ./examples")
}

func processFile(filename string) (*compiler.Service, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	tokens := compiler.Tokenize(string(content))
	parser := compiler.NewParser(tokens)
	service, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("compile error: %w", err)
	}

	return service, nil
}

func generateCode(filename string) string {
	service, err := processFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	outName := baseName + ".go"

	err = compiler.GenerateCode(service, "generated", outName)
	if err != nil {
		fmt.Printf("failed to generate code: %s\n", err)
		os.Exit(1)
	}

	outPath := filepath.Join("generated", outName)
	fmt.Printf("Generated Go code at %s\n", outPath)
	return outPath
}

func buildBinary(filename string) {
	outPath := generateCode(filename)
	baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	binName := baseName
	if os.PathSeparator == '\\' {
		binName += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binName, outPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Building binary:", binName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Build failed: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Build successful.")
}

func runService(filename string) {
	outPath := generateCode(filename)

	cmd := exec.Command("go", "run", outPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running service...")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Service exited with error: %s\n", err)
		os.Exit(1)
	}
}
