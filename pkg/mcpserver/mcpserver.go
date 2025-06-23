package mcpserver

import (
	"context"
	"errors"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	mcpvsphere "github.com/openshift/installer/pkg/mcpserver/vsphere"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/version"
)

type InstallerMcpServer struct {
	Server *server.MCPServer
	Tools  []server.ServerTool
}

// using design examples from https://github.com/Prashanth684/releasecontroller-mcp-server/tree/main

func NewInstallerMcpServer(serverTools []server.ServerTool) *InstallerMcpServer {

	installerMcpServer := &InstallerMcpServer{}

	// todo: yeah I know...

	versionString, _ := version.Version()

	hooks := &server.Hooks{}

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		logrus.Info(session.SessionID())
	})
	s := server.NewMCPServer(
		"OpenShift Installer",
		versionString,
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	installerMcpServer.Server = s
	installerMcpServer.Tools = append(installerMcpServer.Tools, serverTools...)
	installerMcpServer.Tools = append(installerMcpServer.Tools, tools()...)

	s.AddTools(installerMcpServer.Tools...)

	logrus.Infof("There are %d installed tools", len(installerMcpServer.Tools))

	for _, tool := range installerMcpServer.Tools {
		logrus.Info("Installing tool: ", tool.Tool.Name)
	}

	return installerMcpServer
}

func tools() []server.ServerTool {
	logrus.Info("Initializing MCP Server Tools")
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("get_coreos_images", mcp.WithDescription("Gets the coreos images in json from the installer")),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return ProcessResults(getCoreOS()), nil
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
				server, ok := arguments["server"].(string)
				if !ok {
					return nil, errors.New("server is required")
				}

				return ProcessResults(mcpvsphere.GetVSphereTopology(username, password, server)), nil
			},
		},
	}
}

func ProcessResults(content string, err error) *mcp.CallToolResult {
	if err != nil {
		return mcp.NewToolResultErrorFromErr(content, err)
	}
	return mcp.NewToolResultText(content)
}

func (i *InstallerMcpServer) RunServeStdio() error {
	logrus.Info("Starting MCP Server")
	return server.ServeStdio(i.Server)
}

func (i *InstallerMcpServer) RunSSEServer() error {
	sseServer := server.NewSSEServer(i.Server,
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(time.Minute))
	logrus.Info("Starting MCP SSE Server")
	return sseServer.Start(":8080")
}

func getCoreOS() (string, error) {
	logrus.Info("Getting CoreOS stream data")
	streamData, err := rhcos.FetchRawCoreOSStream(context.Background())
	if err != nil {
		return "", err
	}
	return string(streamData), nil
}
