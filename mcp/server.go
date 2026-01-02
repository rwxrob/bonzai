package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rwxrob/bonzai"
)

// config holds configuration options for server creation
type config struct {
	onlyTagged bool
	version    string
}

// Option configures how the MCP server is created
type Option func(*config)

// OnlyTagged returns an Option that configures the server to only
// expose commands that have Mcp metadata set (Mcp != nil).
func OnlyTagged() Option {
	return func(c *config) {
		c.onlyTagged = true
	}
}

// WithVersion sets the server version string.
func WithVersion(version string) Option {
	return func(c *config) {
		c.version = version
	}
}

// NewServer creates an MCP server from a Bonzai command tree.
// It walks the command tree and registers each command as an MCP tool.
// By default, all commands with a Do function are registered.
// Use OnlyTagged() option to only register commands with Mcp metadata.
func NewServer(root *bonzai.Cmd, opts ...Option) *server.MCPServer {
	cfg := &config{
		version: "1.0.0",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	s := server.NewMCPServer(
		root.Name,
		cfg.version,
		server.WithToolCapabilities(true),
	)

	// Walk the command tree and register tools
	root.WalkDeep(func(level int, cmd *bonzai.Cmd) error {
		if !shouldRegister(cmd, cfg) {
			return nil
		}
		registerTool(s, cmd)
		return nil
	}, func(err error) {
		// ignore errors during walk
	})

	return s
}

// shouldRegister determines if a command should be registered as an MCP tool
func shouldRegister(cmd *bonzai.Cmd, cfg *config) bool {
	// Must have a Do function to be callable
	if cmd.Do == nil {
		return false
	}

	// If onlyTagged, require Mcp metadata
	if cfg.onlyTagged && cmd.Mcp == nil {
		return false
	}

	return true
}

// registerTool adds a command as an MCP tool to the server
func registerTool(s *server.MCPServer, cmd *bonzai.Cmd) {
	toolOpts := BuildToolOptions(cmd)
	tool := mcp.NewTool(ToolName(cmd), toolOpts...)
	handler := WrapHandler(cmd)
	s.AddTool(tool, handler)
}
