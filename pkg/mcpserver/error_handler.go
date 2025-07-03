package mcpserver

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// cursor wrote this below
// it doesn't work and probably can be removed

// MCPErrorHandler handles errors by sending them to MCP clients instead of exiting
type MCPErrorHandler struct {
	mcpServer *server.MCPServer
	ctx       context.Context
}

// NewMCPErrorHandler creates a new error handler that sends errors to MCP clients
func NewMCPErrorHandler(mcpServer *server.MCPServer, ctx context.Context) *MCPErrorHandler {
	return &MCPErrorHandler{
		mcpServer: mcpServer,
		ctx:       ctx,
	}
}

// HandleError sends an error notification to all MCP clients instead of exiting
func (h *MCPErrorHandler) HandleError(err error, context string) {
	if h.mcpServer == nil {
		// Fallback to standard error handling if MCP server is not available
		logrus.Errorf("MCP Error Handler: %s - %v", context, err)
		return
	}

	// Send error notification to all clients
	h.mcpServer.SendNotificationToAllClients(
		"notifications/error",
		map[string]any{
			"error":   err.Error(),
			"context": context,
			"type":    "installer_error",
		},
	)

	// Also log the error locally
	logrus.Errorf("MCP Error Handler: %s - %v", context, err)
}

// HandleFatalError handles fatal errors by sending them to MCP clients instead of exiting
func (h *MCPErrorHandler) HandleFatalError(err error, context string) {
	h.HandleError(err, context)

	// In MCP mode, we don't exit - just log and continue
	logrus.Errorf("Fatal error handled by MCP Error Handler: %s - %v", context, err)
}

// HandlePanic recovers from panics and sends them to MCP clients
func (h *MCPErrorHandler) HandlePanic() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		err := fmt.Errorf("panic recovered: %v", r)

		h.HandleError(err, "panic")

		// Log the stack trace
		logrus.Errorf("Panic stack trace:\n%s", stack)
	}
}

// IsMCPMode returns true if the installer is running in MCP mode
func IsMCPMode() bool {
	// Check if any MCP-related environment variables are set
	return os.Getenv("MCP_MODE") == "true" ||
		os.Getenv("MCP_SERVER") == "true" ||
		len(os.Args) > 1 && (os.Args[1] == "run" ||
			os.Args[1] == "mcp-server-sse" ||
			os.Args[1] == "mcp-server-streamable-http" ||
			os.Args[1] == "mcp-server-stdio")
}
