package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	mcpvsphere "github.com/openshift/installer/pkg/mcpserver/vsphere"
	"github.com/openshift/installer/pkg/version"
)

type InstallerMcpServer struct {
	Server    *server.MCPServer
	Tools     []server.ServerTool
	Resources []server.ServerResource
	Prompts   []server.ServerPrompt
}

// using design examples from https://github.com/Prashanth684/releasecontroller-mcp-server/tree/main

func NewInstallerMcpServer(serverTools []server.ServerTool, resources []server.ServerResource) *InstallerMcpServer {
	installerMcpServer := &InstallerMcpServer{}

	versionString, err := version.Version()
	if err != nil {
		logrus.Warn(err)
	}

	hooks := &server.Hooks{}

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		logrus.Info(session.SessionID())
	})
	s := server.NewMCPServer(
		"OpenShift Installer",
		versionString,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	installerMcpServer.Server = s
	installerMcpServer.Tools = append(installerMcpServer.Tools, serverTools...)
	// todo: this could be per platform Tools()
	installerMcpServer.Tools = append(installerMcpServer.Tools, Tools()...)

	installerMcpServer.Resources = append(installerMcpServer.Resources, resources...)

	s.AddTools(installerMcpServer.Tools...)
	s.AddResources(installerMcpServer.Resources...)

	logrus.Infof("There are %d installed tools", len(installerMcpServer.Tools))

	for _, tool := range installerMcpServer.Tools {
		logrus.Info("Installing tool: ", tool.Tool.Name)
	}

	return installerMcpServer
}

func Tools() []server.ServerTool {
	logrus.Info("Initializing MCP Server Tools")
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("get_coreos_images", mcp.WithDescription("Gets the coreos images in json from the installer")),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return ProcessResults(GetCoreOS()), nil
			},
		},
		{
			Tool: mcp.NewTool("get_vsphere_topology",
				mcp.WithDescription("Gets the vsphere topology in json from the installer"),
				mcp.WithString("username"),
				mcp.WithString("password"),
				mcp.WithString("server"),
			),
			Handler: func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				arguments := request.GetArguments()
				// check if the arguments are present using ok pattern
				username, ok := arguments["username"].(string)
				if !ok {
					return nil, errors.New("username is required")
				}
				password, ok := arguments["password"].(string)
				if !ok {
					return nil, errors.New("password is required")
				}
				// missing server
				vServer, ok := arguments["server"].(string)
				if !ok {
					return nil, errors.New("server is required")
				}

				return ProcessResults(
					mcpvsphere.GetVSphereTopology(username, password, vServer),
				), nil
			},
		},
	}
}

func ProcessResourceResults(content string, err error) (*mcp.ReadResourceResult, error) {
	if err != nil {
		return nil, err
	}
	return mcp.NewReadResourceResult(content), nil

}

func ProcessResults(content string, err error) *mcp.CallToolResult {
	if err != nil {
		return mcp.NewToolResultErrorFromErr(content, err)
	}

	/*
		if json.Valid([]byte(content)) {
			return mcp.newtoolre

		}

	*/

	return mcp.NewToolResultText(content)
}

func (i *InstallerMcpServer) RunServeStdio() error {
	logrus.Info("Starting MCP Server")
	return server.ServeStdio(i.Server)
}

func (i *InstallerMcpServer) RunSSEServer() error {
	sseServer := server.NewSSEServer(i.Server,
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second))
	logrus.Info("Starting MCP SSE Server")
	return sseServer.Start(":8080")
}

type McpLogrusHook struct {
	// LogLevels specifies which log levels should trigger this hook.
	LogLevels []logrus.Level

	MCPServer *server.MCPServer

	ProgressToken mcp.ProgressToken
	ClientContext context.Context
}

func (hook *McpLogrusHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		// If formatting fails, we return an error.
		return fmt.Errorf("failed to format log entry: %w", err)
	}
	//logMsgNotify := mcp.NewLoggingMessageNotification(mcp.LoggingLevel(entry.Level.String()), "logrus", line)

	err = hook.MCPServer.SendNotificationToClient(hook.ClientContext,
		"notifications/message",
		map[string]any{
			"level":  entry.Level.String(),
			"data":   line,
			"logger": "logrus",
		},
	)
	if err != nil {
		logrus.Warnf("Failed to send notification to MCP server: %v", err)
	}

	return err
}
func (hook *McpLogrusHook) Levels() []logrus.Level {
	return hook.LogLevels
}
