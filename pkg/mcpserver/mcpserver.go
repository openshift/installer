package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

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

	s := server.NewMCPServer(
		"OpenShift Installer",
		versionString,
		server.WithToolCapabilities(true),
		server.WithLogging(),
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
	sseServer := server.NewSSEServer(i.Server)
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
