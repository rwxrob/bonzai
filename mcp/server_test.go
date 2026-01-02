package mcp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rwxrob/bonzai"
	bmcp "github.com/rwxrob/bonzai/mcp"
)

func TestNewServer_Basic(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Do: func(x *bonzai.Cmd, args ...string) error {
			fmt.Println("executed")
			return nil
		},
	}

	s := bmcp.NewServer(cmd)
	if s == nil {
		t.Fatal("expected server to be created")
	}
}

func TestNewServer_WithMcpMeta(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "greet",
		Short: "greet someone",
		Mcp: &bonzai.McpMeta{
			Desc: "greet a person by name",
			Params: []bonzai.McpParam{
				{Name: "name", Desc: "person's name", Type: "string", Required: true},
			},
		},
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	s := bmcp.NewServer(cmd)
	if s == nil {
		t.Fatal("expected server to be created")
	}
}

func TestNewServer_WithSubcommands(t *testing.T) {
	subCmd := &bonzai.Cmd{
		Name:  "sub",
		Short: "subcommand",
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	rootCmd := &bonzai.Cmd{
		Name:  "root",
		Short: "root command",
		Cmds:  []*bonzai.Cmd{subCmd},
	}

	s := bmcp.NewServer(rootCmd)
	if s == nil {
		t.Fatal("expected server to be created")
	}
}

func TestNewServer_OnlyTagged(t *testing.T) {
	taggedCmd := &bonzai.Cmd{
		Name:  "tagged",
		Short: "tagged command",
		Mcp: &bonzai.McpMeta{
			Desc: "this has mcp metadata",
		},
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	untaggedCmd := &bonzai.Cmd{
		Name:  "untagged",
		Short: "untagged command",
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	rootCmd := &bonzai.Cmd{
		Name:  "root",
		Short: "root command",
		Cmds:  []*bonzai.Cmd{taggedCmd, untaggedCmd},
	}

	s := bmcp.NewServer(rootCmd, bmcp.OnlyTagged())
	if s == nil {
		t.Fatal("expected server to be created")
	}
}

func TestBuildToolOptions_String(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Mcp: &bonzai.McpMeta{
			Params: []bonzai.McpParam{
				{Name: "input", Desc: "input value", Type: "string", Required: true},
			},
		},
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	opts := bmcp.BuildToolOptions(cmd)
	if len(opts) == 0 {
		t.Fatal("expected at least one option")
	}
}

func TestBuildToolOptions_AllTypes(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Mcp: &bonzai.McpMeta{
			Params: []bonzai.McpParam{
				{Name: "str", Desc: "string param", Type: "string", Required: true},
				{Name: "num", Desc: "number param", Type: "number", Required: false},
				{Name: "bool", Desc: "bool param", Type: "boolean", Required: false},
				{Name: "arr", Desc: "array param", Type: "array", Required: false},
			},
		},
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	opts := bmcp.BuildToolOptions(cmd)
	if len(opts) != 5 { // description + 4 params
		t.Fatalf("expected 5 options, got %d", len(opts))
	}
}

func TestBuildToolOptions_WithEnum(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Mcp: &bonzai.McpMeta{
			Params: []bonzai.McpParam{
				{Name: "level", Desc: "log level", Type: "string", Enum: []string{"debug", "info", "warn", "error"}},
			},
		},
		Do: func(x *bonzai.Cmd, args ...string) error {
			return nil
		},
	}

	opts := bmcp.BuildToolOptions(cmd)
	if len(opts) == 0 {
		t.Fatal("expected options for enum param")
	}
}

func TestWrapHandler(t *testing.T) {
	var executed bool

	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Do: func(x *bonzai.Cmd, args ...string) error {
			executed = true
			return nil
		},
	}

	handler := bmcp.WrapHandler(cmd)
	if handler == nil {
		t.Fatal("expected handler to be created")
	}

	// Create a mock request
	req := mcp.CallToolRequest{}
	req.Params.Name = "test"
	req.Params.Arguments = map[string]any{
		"arg1": "value1",
		"arg2": "value2",
	}

	ctx := context.Background()
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !executed {
		t.Fatal("expected command to be executed")
	}

	if result == nil {
		t.Fatal("expected result")
	}
}

func TestWrapHandler_WithError(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Do: func(x *bonzai.Cmd, args ...string) error {
			return fmt.Errorf("test error")
		},
	}

	handler := bmcp.WrapHandler(cmd)
	req := mcp.CallToolRequest{}
	req.Params.Name = "test"

	ctx := context.Background()
	result, err := handler(ctx, req)

	// MCP handlers should return error in result, not as error
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result even on command error")
	}
}

func TestToolName_Simple(t *testing.T) {
	cmd := &bonzai.Cmd{
		Name:  "test",
		Short: "test command",
		Do:    func(x *bonzai.Cmd, args ...string) error { return nil },
	}

	name := bmcp.ToolName(cmd)
	if name != "test" {
		t.Fatalf("expected 'test', got '%s'", name)
	}
}

func TestToolName_Nested(t *testing.T) {
	sub := &bonzai.Cmd{
		Name:  "sub",
		Short: "sub command",
		Do:    func(x *bonzai.Cmd, args ...string) error { return nil },
	}

	root := &bonzai.Cmd{
		Name:  "root",
		Short: "root command",
		Cmds:  []*bonzai.Cmd{sub},
	}

	// Simulate calling Seek to set up caller chain
	root.Seek("sub")

	name := bmcp.ToolName(sub)
	// For a simple command without caller set, it should just return the name
	if name != "sub" {
		t.Fatalf("expected 'sub', got '%s'", name)
	}
}
