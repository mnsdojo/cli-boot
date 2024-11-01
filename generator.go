package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// File structure definitions
var fileContents = map[string]string{
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

// Add your package functionality here.
`,
	"go.mod": `module {{.ModuleName}}

go {{.GoVersion}}
`,
	"main.go": `package main


import "github.com/{{.ModuleName}}/pkg"

func main() {
    // Entry point for your application
}
`,

	"README.md": `# {{.ProjectName}}
This is a simple CLI application created with Go.
`,
}

func displayAsciiArt() {
	asciiArt := `
  

   ____ _     ___ ____   ___   ___ _____ 
  / ___| |   |_ _| __ ) / _ \ / _ \_   _|
 | |   | |    | ||  _ \| | | | | | || |  
 | |___| |___ | || |_) | |_| | |_| || |  
  \____|_____|___|____/ \___/ \___/ |_|  
                                         


    `
	fmt.Println(asciiArt)
}

func prompt(msg string) string {
	fmt.Println(msg, "")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	displayAsciiArt()
	projectName := prompt("Enter the project name:")
	moduleName := prompt("Enter the module name: ")
	goVersion := prompt("Enter the go version (default is 1.23.1)")
	if goVersion == "" {
		goVersion = "1.23.1"
	}

	createProjectStructure(projectName, moduleName, goVersion)
}

func createProjectStructure(projectName, moduleName, goVersion string) {
	if err := os.Mkdir(projectName, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating project dirrectory : %v\n", err)
		return
	}

	// paths structure
	paths := []string{
		filepath.Join(projectName, "cmd"),
		filepath.Join(projectName, "pkg"),
		filepath.Join(projectName, "cmd", "root.go"),
		filepath.Join(projectName, "pkg", "mypackage.go"),
		filepath.Join(projectName, "go.mod"),
		filepath.Join(projectName, "main.go"),
		filepath.Join(projectName, "README.md"),
	}

	for _, path := range paths {
		if err := createPath(path, projectName, moduleName, goVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating : %s : %v\n", path, err)
		}
		fmt.Printf("Project structure created : %s\n", projectName)
	}
}

// createPath creates directories or files
func createPath(path, projectName, moduleName, goVersion string) error {
	if filepath.Ext(path) == "" {
		// Create directories
		return os.MkdirAll(path, 0755)
	}
	// Create files with content
	return createFile(path, projectName, moduleName, goVersion)
}

// createFile writes content to a file, replacing placeholders with project-specific values
func createFile(path, projectName, moduleName, goVersion string) error {
	content, exists := fileContents[filepath.Base(path)]
	if !exists {
		return fmt.Errorf("no predefined content for %s", path)
	}
	content = replacePlaceholders(content, projectName, moduleName, goVersion)
	return os.WriteFile(path, []byte(content), 0644)
}

// replacePlaceholders replaces placeholders in the template
func replacePlaceholders(content, projectName, moduleName, goVersion string) string {
	content = replace(content, "{{.ProjectName}}", projectName)
	content = replace(content, "{{.ModuleName}}", moduleName)
	content = replace(content, "{{.GoVersion}}", goVersion)
	return content
}

func replace(content, placeholder, value string) string {
	return strings.ReplaceAll(content, placeholder, value)
}
