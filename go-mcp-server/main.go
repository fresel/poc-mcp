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
	ProjectName string `json:"projectName" jsonschema:"name of the Java project to create or extend"`
	Folder      string `json:"folder" jsonschema:"folder relative to the workspace root where the project should be created"`
	Kind        string `json:"kind" jsonschema:"template kind to add: controller, service, integration-test, or unit-test"`
}

type ScaffoldOutput struct {
	Path string `json:"path" jsonschema:"path of the created Java application"`
}

func ScaffoldJavaApp(ctx context.Context, req *mcp.CallToolRequest, input ScaffoldInput) (
	*mcp.CallToolResult,
	ScaffoldOutput,
	error,
) {
	if strings.TrimSpace(input.ProjectName) == "" {
		return nil, ScaffoldOutput{}, fmt.Errorf("projectName is required")
	}
	if strings.TrimSpace(input.Folder) == "" {
		return nil, ScaffoldOutput{}, fmt.Errorf("folder is required")
	}
	if !validScaffoldKinds[input.Kind] {
		return nil, ScaffoldOutput{}, fmt.Errorf("kind must be one of: controller, service, model, integration-test, unit-test")
	}

	workspaceRoot := os.Getenv("MCP_WORKSPACE_ROOT")
	if workspaceRoot == "" {
		workspaceRoot = ".."
	}
	workspaceRoot, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return nil, ScaffoldOutput{}, fmt.Errorf("resolve workspace root: %w", err)
	}

	folder := filepath.Clean(input.Folder)
	if filepath.IsAbs(folder) {
		return nil, ScaffoldOutput{}, fmt.Errorf("folder must be relative to the workspace root")
	}
	parent := filepath.Join(workspaceRoot, folder)
	projectName := filepath.Clean(input.ProjectName)
	if projectName == "." || projectName == ".." || filepath.IsAbs(projectName) || strings.Contains(projectName, string(filepath.Separator)) {
		return nil, ScaffoldOutput{}, fmt.Errorf("projectName must be a single directory name")
	}
	target := filepath.Join(parent, projectName)
	relative, err := filepath.Rel(workspaceRoot, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil, ScaffoldOutput{}, fmt.Errorf("project path must stay inside the workspace root")
	}
	_, statErr := os.Stat(target)
	projectExists := statErr == nil
	if statErr != nil && !os.IsNotExist(statErr) {
		return nil, ScaffoldOutput{}, fmt.Errorf("check project: %w", statErr)
	}
	log.Printf("scaffold_java_app: project=%s kind=%s target=%s", input.ProjectName, input.Kind, target)

	templateRoot := filepath.Join("templates", "base")
	kindRoot := filepath.Join("templates", input.Kind)
	templateFiles := []string{}
	if !projectExists {
		templateFiles = append(templateFiles, filepath.Join(templateRoot, "pom.xml"))
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
		if _, err := os.Stat(path); err == nil {
			return nil, ScaffoldOutput{}, fmt.Errorf("scaffold file already exists: %s", path)
		} else if !os.IsNotExist(err) {
			return nil, ScaffoldOutput{}, fmt.Errorf("check scaffold file: %w", err)
		}
		if name == "pom.xml" {
			content = []byte(strings.Replace(string(content), "<artifactId>java-test</artifactId>", "<artifactId>"+input.ProjectName+"</artifactId>", 1))
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, ScaffoldOutput{}, fmt.Errorf("create directory for %s: %w", name, err)
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return nil, ScaffoldOutput{}, fmt.Errorf("write %s: %w", name, err)
		}
	}

	log.Printf("scaffold_java_app: created %s", target)
	return nil, ScaffoldOutput{Path: target}, nil
}

var validScaffoldKinds = map[string]bool{
	"controller":       true,
	"service":          true,
	"model":            true,
	"integration-test": true,
	"unit-test":        true,
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "scaffold_java_app",
		Description: "create or extend a named Java Maven app under a workspace-relative folder",
	}, ScaffoldJavaApp)
	// Run the server over stdin/stdout, until the client disconnects.
	log.Printf("starting greeter MCP server")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
