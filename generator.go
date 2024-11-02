package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

type ProjectConfig struct {
	ProjectName string
	ModuleName  string
	GoVersion   string
}

type FileTemplate struct {
	Path     string
	Template string
}

var fileTemplates = map[string]string{
	"root.go": `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from {{.ProjectName}}!")
	os.Exit(0)
}`,
	"mypackage.go": `package mypackage

// Package mypackage provides core functionality for {{.ProjectName}}
`,
	"go.mod": `module {{.ModuleName}}

go {{.GoVersion}}
`,
	"main.go": `package main

import (
	"log"

	"github.com/{{.ModuleName}}/pkg"
)

func main() {
	// Entry point for your application
	log.Println("Starting {{.ProjectName}}...")
}
`,
	"README.md": `# {{.ProjectName}}

## Overview
This is a CLI application created with Go.

## Installation
` + "```bash" + `
go get github.com/{{.ModuleName}}
` + "```" + `

## Usage
Describe how to use your application here.

## License
MIT License
`,
}

var (
	success = color.New(color.FgGreen).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
	warn    = color.New(color.FgYellow).SprintFunc()
	fail    = color.New(color.FgRed).SprintFunc()
)

func displayAsciiArt() {
	asciiArt := `
  

   ____ _     ___ ____   ___   ___ _____ 
  / ___| |   |_ _| __ ) / _ \ / _ \_   _|
 | |   | |    | ||  _ \| | | | | | || |  
 | |___| |___ | || |_) | |_| | |_| || |  
  \____|_____|___|____/ \___/ \___/ |_|  
                                         


    `
	fmt.Println(info(asciiArt))
}

func prompt(msg string, defaultVal string) string {
	if defaultVal != " " {
		fmt.Printf("%s (default : %s) : ", msg, defaultVal)
	} else {
		fmt.Printf("%s", msg)
	}
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input :%v", err)
	}
	input = strings.TrimSpace(input)
	if input == " " && defaultVal != "" {
		return defaultVal
	}
	return input
}

func validateInput(config *ProjectConfig) error {
	if config.ProjectName == " " {
		return fmt.Errorf("project name cannot be empty")
	}
	if config.ModuleName == "" {
		return fmt.Errorf("module cannot be empty")
	}
	if !strings.HasPrefix(config.GoVersion, "1.") {
		return fmt.Errorf("invalid go version format")
	}

	return nil
}

func createProjectStructure(config *ProjectConfig) error {
	if err := os.Mkdir(config.ProjectName, 0755); err != nil {
		return fmt.Errorf("error creating project directory : %w", err)
	}

	// First create list of dirs
	// paths structure
	dirs := []string{
		filepath.Join(config.ProjectName, "cmd"),
		filepath.Join(config.ProjectName, "pkg"),
		filepath.Join(config.ProjectName, "internal"),
		filepath.Join(config.ProjectName, "docs"),
		filepath.Join(config.ProjectName, "scripts"),
	}

	// create dirs
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directories %s :%w", dir, err)
		}
		fmt.Printf("%s Created directory: %s\n", success("✓"), dir)

	}
	// Create files from templates --
	files := map[string]string{
		filepath.Join(config.ProjectName, "cmd", "root.go"):      fileTemplates["root.go"],
		filepath.Join(config.ProjectName, "pkg", "mypackage.go"): fileTemplates["mypackage.go"],
		filepath.Join(config.ProjectName, "go.mod"):              fileTemplates["go.mod"],
		filepath.Join(config.ProjectName, "main.go"):             fileTemplates["main.go"],
		filepath.Join(config.ProjectName, "README.md"):           fileTemplates["README.md"],
	}

	for path, content := range files {
		if err := createFileFromTemplate(path, content, config); err != nil {
			return fmt.Errorf("failed to parse template :%w", err)
		}
		fmt.Printf("%s Created file: %s\n", success("✓"), path)

	}
	return nil
}

func createFileFromTemplate(path, content string, config *ProjectConfig) error {
	tmpl, err := template.New(filepath.Base(path)).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file : %w", err)
	}
	defer file.Close()
	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}

func main() {
	log.SetFlags(0)
	displayAsciiArt()

	config := &ProjectConfig{
		ProjectName: prompt("Enter the project name", ""),
		ModuleName:  prompt("Enter the module name (e.g., github.com/username/project)", ""),
		GoVersion:   prompt("Enter the Go version", "1.21"),
	}
	if err := validateInput(config); err != nil {
		log.Fatalf("%s Validation error: %v", fail("✗"), err)
	}
	fmt.Printf("\n%s Creating project structure ..\n\n", info("->"))
	if err := createProjectStructure(config); err != nil {
		log.Fatalf("%s Error :%v", fail("x"), err)
	}
	fmt.Printf("\n%s Project %s created successfully!\n", success("✓"), config.ProjectName)
	fmt.Printf("\n%s Next steps : \n", info("->"))
	fmt.Printf("  1. cd %s\n", config.ProjectName)
	fmt.Printf("  2. go mod tidy\n")
	fmt.Printf("  3. go run main.go\n\n")
}
