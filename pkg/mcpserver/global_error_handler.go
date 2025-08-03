package mcpserver

import (
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// cursor wrote this below
// it doesn't work and probably can be removed

var (
	globalErrorHandler *MCPErrorHandler
	globalErrorMutex   sync.RWMutex
)

// SetGlobalErrorHandler sets the global error handler for the installer
func SetGlobalErrorHandler(handler *MCPErrorHandler) {
	globalErrorMutex.Lock()
	defer globalErrorMutex.Unlock()
	globalErrorHandler = handler
}

// GetGlobalErrorHandler returns the global error handler
func GetGlobalErrorHandler() *MCPErrorHandler {
	globalErrorMutex.RLock()
	defer globalErrorMutex.RUnlock()
	return globalErrorHandler
}

// HandleGlobalError handles errors globally, sending them to MCP clients if in MCP mode
func HandleGlobalError(err error, context string) {
	if IsMCPMode() {
		handler := GetGlobalErrorHandler()
		if handler != nil {
			handler.HandleFatalError(err, context)
			return
		}
	}

	// Fallback to standard error handling
	logrus.Errorf("Global Error Handler: %s - %v", context, err)
}

// HandleGlobalFatal handles fatal errors globally, preventing exit in MCP mode
func HandleGlobalFatal(err error, context string) {
	if IsMCPMode() {
		handler := GetGlobalErrorHandler()
		if handler != nil {
			handler.HandleFatalError(err, context)
			return
		}
	}

	// In non-MCP mode, exit as usual
	logrus.Fatalf("Fatal error: %s - %v", context, err)
}

// SafeFatal is a replacement for logrus.Fatal that respects MCP mode
func SafeFatal(err error, context string) {
	HandleGlobalFatal(err, context)
}

// SafeFatalf is a replacement for logrus.Fatalf that respects MCP mode
func SafeFatalf(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	HandleGlobalFatal(err, "formatted error")
}

// SafeExit is a replacement for logrus.Exit that respects MCP mode
func SafeExit(code int, context string) {
	if IsMCPMode() {
		handler := GetGlobalErrorHandler()
		if handler != nil {
			err := fmt.Errorf("exit called with code %d: %s", code, context)
			handler.HandleFatalError(err, "exit called")
			return
		}
	}

	// In non-MCP mode, exit as usual
	logrus.Exit(code)
}

// SetupMCPErrorHandling sets up global error handling for MCP mode
func SetupMCPErrorHandling(mcpServer *server.MCPServer, ctx context.Context) {
	if IsMCPMode() {
		handler := NewMCPErrorHandler(mcpServer, ctx)
		SetGlobalErrorHandler(handler)

		logrus.Info("MCP error handling enabled - fatal errors will be sent to MCP clients instead of exiting")
	}
}
