package main

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

func main() {
}
