package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/command"
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
	cmd.AddCommand(newRunMcpStreamableHttpCmd())

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

			i := mcpserver.NewInstallerMcpServer(Tools(), Resources(), ResourceTemplates())

			err := i.RunSSEServer()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

// https://mcp-go.dev/transports/http#standard-mcp-endpoints
func newRunMcpStreamableHttpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp-server-streamable-http",
		Short: "Run MCP Server-Sent Events",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			i := mcpserver.NewInstallerMcpServer(Tools(), Resources(), ResourceTemplates())

			err := i.RunStreamableHttp()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newRunMcpStdioCmd() *cobra.Command {

	// this wouldn't work as-is since logging is sent to stdout
	return &cobra.Command{
		Use:   "mcp-server-stdio",
		Short: "Run MCP Server Stdio",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()
			logrus.Fatal("not implemented")
		},
	}
}
