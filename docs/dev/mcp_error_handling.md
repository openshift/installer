# MCP Error Handling

## Overview

When running the OpenShift installer as an MCP (Model Context Protocol) server, the installer should never exit on errors. Instead, errors should be sent to MCP clients as notifications. This document describes the error handling approach implemented for MCP mode.

## Key Components

### 1. MCPErrorHandler

The `MCPErrorHandler` is responsible for:
- Sending error notifications to all connected MCP clients
- Logging errors locally for debugging
- Preventing the installer from exiting on fatal errors

### 2. Global Error Handler

The global error handler provides:
- Centralized error handling for the entire installer
- Safe replacements for `logrus.Fatal`, `logrus.Fatalf`, and `logrus.Exit`
- Automatic detection of MCP mode

### 3. MCP Mode Detection

MCP mode is detected by:
- Environment variables: `MCP_MODE=true` or `MCP_SERVER=true`
- Command line arguments: `run`, `mcp-server-sse`, `mcp-server-streamable-http`, `mcp-server-stdio`

## Usage

### For MCP Server Commands

```go
// Set up global error handling
ctx := context.Background()
mcpserver.SetupMCPErrorHandling(mcpServer, ctx)

// Use safe error handling functions
err := someOperation()
if err != nil {
    mcpserver.SafeFatal(err, "operation context")
}
```

### For Other Commands

Replace direct calls to `logrus.Fatal` with `mcpserver.SafeFatal`:

```go
// Instead of:
logrus.Fatal(err)

// Use:
mcpserver.SafeFatal(err, "context description")
```

### Error Notifications

Errors are sent to MCP clients as notifications with the following structure:

```json
{
  "jsonrpc": "2.0",
  "method": "notifications/error",
  "params": {
    "error": "error message",
    "context": "context description",
    "type": "installer_error"
  }
}
```

## Implementation Details

### Error Handler Setup

1. **MCP Server Creation**: When creating an MCP server, error handling hooks are automatically added
2. **Global Handler**: The global error handler is set up when `SetupMCPErrorHandling` is called
3. **Mode Detection**: The system automatically detects if it's running in MCP mode

### Error Flow

1. **Error Occurs**: An error occurs somewhere in the installer
2. **Handler Intercepts**: The global error handler intercepts the error
3. **MCP Mode Check**: If in MCP mode, the error is sent to all connected clients
4. **No Exit**: The installer continues running instead of exiting
5. **Logging**: The error is logged locally for debugging

### Fallback Behavior

If MCP mode is not detected or the error handler is not available:
- Standard error handling is used
- The installer may exit as usual
- Errors are logged normally

## Migration Guide

To migrate existing code to use MCP error handling:

1. **Replace Fatal Calls**:
   ```go
   // Old
   logrus.Fatal(err)
   
   // New
   mcpserver.SafeFatal(err, "context")
   ```

2. **Replace Fatalf Calls**:
   ```go
   // Old
   logrus.Fatalf("Error: %v", err)
   
   // New
   mcpserver.SafeFatalf("Error: %v", err)
   ```

3. **Replace Exit Calls**:
   ```go
   // Old
   logrus.Exit(code)
   
   // New
   mcpserver.SafeExit(code, "context")
   ```

## Testing

To test MCP error handling:

1. **Set Environment Variable**:
   ```bash
   export MCP_MODE=true
   ```

2. **Run MCP Server**:
   ```bash
   openshift-install run mcp-server-sse
   ```

3. **Trigger Error**: Errors should be sent to MCP clients instead of causing the server to exit

## Future Enhancements

- **Error Categories**: Add support for categorizing errors (warning, error, fatal)
- **Error Recovery**: Implement automatic error recovery mechanisms
- **Client-Specific Notifications**: Send different error notifications to different clients
- **Error Metrics**: Track error statistics and send them to monitoring systems 