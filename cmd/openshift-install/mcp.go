package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/command"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/mcpserver"
)

func newRunCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run MCP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newRunMcpStdioCmd())
	cmd.AddCommand(newRunMcpSseCmd())

	return cmd
}

func newRunMcpSseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp-server-sse",
		Short: "Run MCP Server-Sent Events",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			i := mcpserver.NewInstallerMcpServer(tools())

			err := i.RunSSEServer()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newRunMcpStdioCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp-server-stdio",
		Short: "Run MCP Server Stdio",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()
			i := mcpserver.NewInstallerMcpServer(tools())

			err := i.RunServeStdio()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}
func tools() []server.ServerTool {
	logrus.Info("Initializing MCP Server Tools")
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("get_graph", mcp.WithDescription("Gets the execution graph from the installer")),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcpserver.ProcessResults(RunGraph()), nil
			},
		},
		{
			Tool: mcp.NewTool("create_cluster",
				mcp.WithDescription("Create OpenShift cluster"),
			),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcpserver.ProcessResults(runCreateCluster()), nil
			},
		},
	}
}

func runCreateCluster() (string, error) {
	ctx := context.Background()

	// how do I create install config????
	// or how could use the tui to my advantage ?

	runCommand := runTargetCmd(ctx, targetassets.Cluster...)

	runCommand(clusterTarget.command, []string{})

	exitCode, err := clusterCreatePostRun(ctx)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("exit code %d", exitCode)
	}

	return "", nil
}
