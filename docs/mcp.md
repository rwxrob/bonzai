# Bonzai MCP (Model Context Protocol) Support

Bonzai provides native support for exposing CLI commands as MCP tools, enabling AI assistants like Claude to interact with your command-line applications.

## Overview

The `bonzai/mcp` package automatically generates an MCP server from your Bonzai command tree. Commands with `Mcp` metadata are exposed as tools that AI assistants can call.

## Quick Start

### 1. Add MCP Metadata to Commands

```go
import "github.com/rwxrob/bonzai"

var MyCmd = &bonzai.Cmd{
    Name:  "greet",
    Short: "send a greeting",
    
    // MCP metadata makes this command available as an MCP tool
    Mcp: &bonzai.McpMeta{
        Desc: "Send a personalized greeting message",
        Params: []bonzai.McpParam{
            {Name: "name", Desc: "Name to greet", Type: "string", Required: true},
            {Name: "formal", Desc: "Use formal greeting", Type: "boolean"},
        },
    },
    
    Do: func(x *bonzai.Cmd, args ...string) error {
        // Command implementation
        return nil
    },
}
```

### 2. Create MCP Server Binary

```go
// cmd/myapp-mcp/main.go
package main

import (
    "log"
    "github.com/mark3labs/mcp-go/server"
    "github.com/rwxrob/bonzai/mcp"
    "myapp/cmd"
)

func main() {
    // Create MCP server from command tree
    // OnlyTagged() exposes only commands with Mcp metadata
    s := mcp.NewServer(cmd.Root, mcp.OnlyTagged())
    
    if err := server.ServeStdio(s); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

### 3. Register with Claude Code

```bash
# Build the MCP server
go build -o myapp-mcp ./cmd/myapp-mcp

# Add to Claude Code
claude mcp add --scope user --transport stdio myapp -- ./myapp-mcp
```

## MCP Metadata Types

### McpMeta

```go
type McpMeta struct {
    Desc      string      // Tool description (falls back to Cmd.Short)
    Params    []McpParam  // Parameter definitions
    Examples  []string    // Usage examples
    Streaming bool        // Supports streaming responses
    Resource  string      // MCP resource URI (if applicable)
}
```

### McpParam

```go
type McpParam struct {
    Name     string   // Parameter name
    Desc     string   // Description
    Type     string   // "string", "number", "boolean", "array"
    Required bool     // Is required
    Enum     []string // Allowed values (optional)
    Pattern  string   // Regex pattern (optional)
}
```

## Server Options

### OnlyTagged()

Only expose commands that have `Mcp` metadata set:

```go
s := mcp.NewServer(cmd.Root, mcp.OnlyTagged())
```

Without this option, all commands are exposed with auto-generated schemas.

## Persistent State

For MCP servers that need state across calls (since each call is a fresh process), use bonzai's persisters:

```go
import "github.com/rwxrob/bonzai/persisters/injson"

// Store state in ~/.local/state/myapp/state.json
var state = injson.NewUserState("myapp", "state.json")

// Use in commands
state.Get("key")
state.Set("key", "value")
```

## Examples

### lazywal - Video Wallpaper Manager

[lazywal](https://github.com/BuddhiLW/lazywal) uses bonzai MCP to let AI assistants control desktop wallpapers:

```bash
# Install
go install github.com/BuddhiLW/lazywal/cmd/lazywal-mcp@latest

# Add to Claude
claude mcp add --scope user --transport stdio lazywal -- lazywal-mcp
```

Then ask Claude: *"Set my wallpaper to ~/Videos/ocean.mp4"*

## AI Assistant Configuration

### Claude Code (CLI)

```bash
claude mcp add --scope user --transport stdio <name> -- /path/to/binary
```

### Claude Desktop

Add to `~/.claude.json`:

```json
{
  "mcpServers": {
    "myapp": {
      "command": "/path/to/myapp-mcp",
      "transport": "stdio"
    }
  }
}
```

### Other MCP-Compatible Assistants

Any assistant supporting the [Model Context Protocol](https://modelcontextprotocol.io/) can use bonzai MCP servers via stdio transport.

## Best Practices

1. **Use descriptive `Desc`**: AI assistants use this to understand when to call your tool
2. **Document parameters**: Clear `McpParam.Desc` helps AI provide correct values
3. **Use `OnlyTagged()`**: Explicitly choose which commands to expose
4. **Handle errors gracefully**: Return clear error messages AI can relay to users
5. **Use persistent state**: For stateful operations, use `injson` or other persisters
