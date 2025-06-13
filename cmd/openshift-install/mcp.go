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

	// this wouldn't work as-is since logging is sent to stdout
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
			Handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcpserver.ProcessResults(runCreateCluster(ctx, req)), nil
			},
		},
	}
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

func runCreateCluster(ctx context.Context, req mcp.CallToolRequest) (string, error) {
	logrus.Info("MCP Server Creating OpenShift cluster")
	srv := server.ServerFromContext(ctx)
	hook := &McpLogrusHook{
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		},
		MCPServer: srv,
	}

	hook.ClientContext = ctx
	progressToken := req.Params.Meta.ProgressToken
	if progressToken != nil {
		hook.ProgressToken = progressToken
	}

	logrus.AddHook(hook)

	// how do I create install config????
	// or how could use the tui to my advantage ?

	i := 0
	steps := 2

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      i,
				"total":         steps,
				"progressToken": progressToken,
				"message":       fmt.Sprintf("Server progress %d", i/steps),
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	runCommand := runTargetCmd(ctx, targetassets.Cluster...)
	runCommand(clusterTarget.command, []string{})

	i = 1

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      i,
				"total":         steps,
				"progressToken": progressToken,
				"message":       fmt.Sprintf("Server progress %d", i/steps),
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}
	exitCode, err := clusterCreatePostRun(ctx)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("exit code %d", exitCode)
	}

	return "", nil
}
