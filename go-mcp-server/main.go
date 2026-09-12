package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed templates
var scaffoldTemplates embed.FS

type Input struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{Greeting: "Hi " + input.Name}, nil
}

type ScaffoldInput struct {
	Folder string `json:"folder" jsonschema:"parent folder where the java-test project should be created; absolute or relative to the server working directory"`
	Kind   string `json:"kind" jsonschema:"template kind: controller, service, integration-test, or unit-test"`
}

type ScaffoldOutput struct {
	Path string `json:"path" jsonschema:"path of the created Java application"`
}

func ScaffoldJavaApp(ctx context.Context, req *mcp.CallToolRequest, input ScaffoldInput) (
	*mcp.CallToolResult,
	ScaffoldOutput,
	error,
) {
	if strings.TrimSpace(input.Folder) == "" {
		return nil, ScaffoldOutput{}, fmt.Errorf("folder is required")
	}
	if !validScaffoldKinds[input.Kind] {
		return nil, ScaffoldOutput{}, fmt.Errorf("kind must be one of: controller, service, integration-test, unit-test")
	}

	parent, err := filepath.Abs(filepath.Clean(input.Folder))
	if err != nil {
		return nil, ScaffoldOutput{}, fmt.Errorf("resolve folder: %w", err)
	}

	target := filepath.Join(parent, "java-test")

	if _, err := os.Stat(target); err == nil {
		return nil, ScaffoldOutput{}, fmt.Errorf("target already exists: %s", target)
	} else if !os.IsNotExist(err) {
		return nil, ScaffoldOutput{}, fmt.Errorf("check target: %w", err)
	}

	templateRoot := filepath.Join("templates", "base")
	kindRoot := filepath.Join("templates", input.Kind)
	templateFiles := []string{
		filepath.Join(templateRoot, "pom.xml"),
	}

	err = fs.WalkDir(scaffoldTemplates, kindRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, ScaffoldOutput{}, fmt.Errorf("load %s templates: %w", input.Kind, err)
	}

	for _, templatePath := range templateFiles {
		content, err := fs.ReadFile(scaffoldTemplates, templatePath)
		if err != nil {
			return nil, ScaffoldOutput{}, fmt.Errorf("read template %s: %w", templatePath, err)
		}

		name := templatePath
		name = strings.TrimPrefix(name, templateRoot+string(filepath.Separator))
		name = strings.TrimPrefix(name, kindRoot+string(filepath.Separator))
		path := filepath.Join(target, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, ScaffoldOutput{}, fmt.Errorf("create directory for %s: %w", name, err)
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return nil, ScaffoldOutput{}, fmt.Errorf("write %s: %w", name, err)
		}
	}

	return nil, ScaffoldOutput{Path: target}, nil
}

var validScaffoldKinds = map[string]bool{
	"controller":       true,
	"service":          true,
	"integration-test": true,
	"unit-test":        true,
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "scaffold_java_app",
		Description: "create a basic Java Maven app named java-test under a user-specified folder",
	}, ScaffoldJavaApp)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
