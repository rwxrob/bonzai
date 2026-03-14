// Package mcp provides auto-generation of MCP (Model Context Protocol)
// servers from Bonzai command trees.
package mcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/BuddhiLW/bonzai"
)

// WrapHandler converts a Cmd.Do function to an MCP handler function.
// It captures stdout and returns it as the tool result content.
func WrapHandler(cmd *bonzai.Cmd) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments from request using the helper method
		args := extractArgs(req.GetArguments())

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the command
		var cmdErr error
		if cmd.Do != nil {
			cmdErr = cmd.Do(cmd, args...)
		}

		// Restore stdout and read captured output
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// If command returned an error, include it in the result
		if cmdErr != nil {
			return mcp.NewToolResultError(cmdErr.Error()), nil
		}

		return mcp.NewToolResultText(output), nil
	}
}

// extractArgs converts MCP arguments map to a slice of strings.
// Arguments are sorted by key for deterministic ordering.
func extractArgs(arguments map[string]any) []string {
	if len(arguments) == 0 {
		return nil
	}

	// Sort keys for deterministic ordering
	keys := make([]string, 0, len(arguments))
	for k := range arguments {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	args := make([]string, 0, len(arguments))
	for _, k := range keys {
		v := arguments[k]
		switch val := v.(type) {
		case string:
			args = append(args, val)
		case []interface{}:
			for _, item := range val {
				args = append(args, fmt.Sprintf("%v", item))
			}
		default:
			args = append(args, fmt.Sprintf("%v", val))
		}
	}
	return args
}

// BuildToolOptions generates mcp.ToolOption slice from Cmd's McpMeta.Params.
// It includes description and all parameter definitions.
func BuildToolOptions(cmd *bonzai.Cmd) []mcp.ToolOption {
	opts := []mcp.ToolOption{}

	// Add description
	desc := cmd.McpDesc()
	if desc != "" {
		opts = append(opts, mcp.WithDescription(desc))
	}

	// If no MCP metadata, return just description
	if cmd.Mcp == nil {
		return opts
	}

	// Add parameter options
	for _, param := range cmd.Mcp.Params {
		opts = append(opts, buildParamOption(param))
	}

	return opts
}

// buildParamOption creates an MCP tool option for a single parameter
func buildParamOption(param bonzai.McpParam) mcp.ToolOption {
	var paramOpts []mcp.PropertyOption

	// Add description if present
	if param.Desc != "" {
		paramOpts = append(paramOpts, mcp.Description(param.Desc))
	}

	// Add required if set
	if param.Required {
		paramOpts = append(paramOpts, mcp.Required())
	}

	// Add enum if present
	if len(param.Enum) > 0 {
		paramOpts = append(paramOpts, mcp.Enum(param.Enum...))
	}

	// Add pattern if present (for string types)
	if param.Pattern != "" {
		paramOpts = append(paramOpts, mcp.Pattern(param.Pattern))
	}

	// Return appropriate type
	switch param.Type {
	case "number":
		return mcp.WithNumber(param.Name, paramOpts...)
	case "boolean":
		return mcp.WithBoolean(param.Name, paramOpts...)
	case "array":
		return mcp.WithArray(param.Name, paramOpts...)
	default: // "string" or unspecified
		return mcp.WithString(param.Name, paramOpts...)
	}
}

// ToolName returns the MCP tool name for a command.
// For simple commands, it returns the command name.
// For nested commands, it returns a path-based name using underscores.
func ToolName(cmd *bonzai.Cmd) string {
	return cmd.Name
}

// ToolNameNested returns a path-based tool name using the command's path.
// Commands are joined with underscores (e.g., "root_sub_leaf").
func ToolNameNested(cmd *bonzai.Cmd) string {
	path := cmd.PathNames()
	if len(path) == 0 {
		return cmd.Name
	}

	// Join with underscores
	name := ""
	for i, p := range path {
		if i > 0 {
			name += "_"
		}
		name += p
	}
	if name == "" {
		return cmd.Name
	}
	return name
}
