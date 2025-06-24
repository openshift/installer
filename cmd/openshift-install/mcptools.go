package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset/cluster"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/mcpserver"
)

func runCreateCluster(ctx context.Context, req mcp.CallToolRequest) (string, error) {
	logrus.Info("MCP Server Creating OpenShift cluster")
	srv := server.ServerFromContext(ctx)
	hook := &mcpserver.McpLogrusHook{
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

	// Set the install directory
	cluster.InstallDir = command.RootOpts.Dir

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := srv.SendNotificationToClient(
					ctx,
					"notifications/message",
					map[string]any{
						"level":   "info",
						"message": "Still working...",
						"logger":  "installer",
					},
				)
				if err != nil {
					logrus.Warn(err)
				}
			case <-done:
				return
			}
		}
	}()
	// Create install-config first
	installConfigRunner := runTargetCmd(ctx, targetassets.InstallConfig...)
	installConfigRunner(installConfigTarget.command, []string{})

	// Now create the cluster
	runCommand := runTargetCmd(ctx, targetassets.Cluster...)
	runCommand(clusterTarget.command, []string{})

	exitCode, err := clusterCreatePostRun(ctx)

	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("exit code %d", exitCode)
	}

	close(done)
	return "Cluster created successfully", nil
}
func Tools() []server.ServerTool {
	logrus.Info("Initializing MCP Server Tools")
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("get_graph", mcp.WithDescription("Gets the execution graph from the installer")),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcpserver.ProcessResults(RunGraph()), nil
			},
		},
		/*
			{
				Tool: mcp.NewTool("create_install_config",
					mcp.WithDescription("Create OpenShift install-config"),
				),
				Handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					return mcpserver.ProcessResults(runCreateInstallConfig(ctx, req)), nil
				},
			},

		*/
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
