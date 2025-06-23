package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset/cluster"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/mcpserver"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/defaults"
	"github.com/openshift/installer/pkg/types/gcp"
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
			Tool: mcp.NewTool("create_install_config",
				mcp.WithDescription("Create OpenShift install-config"),
			),
			Handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcpserver.ProcessResults(runCreateInstallConfig(ctx, req)), nil
			},
		},
		{
			Tool: mcp.NewTool("create_custom_install_config",
				mcp.WithDescription("Create OpenShift install-config with custom parameters"),
			),
			Handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcpserver.ProcessResults(runCreateCustomInstallConfig(ctx, req)), nil
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

// createInstallConfig creates a basic install-config for testing
func createInstallConfig() *types.InstallConfig {
	clusterNetwork, _ := ipnet.ParseCIDR("10.128.0.0/14")
	serviceNetwork, _ := ipnet.ParseCIDR("172.30.0.0/16")
	machineNetwork, _ := ipnet.ParseCIDR("10.0.0.0/16")

	config := &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		BaseDomain: "example.com",
		SSHKey:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDK6UTEydcEKzuNdPaofn8Z2DwgHqdcionLZBiPf/zIRNco++etLsat7Avv7yt04DINQd5zjxIFgG8jblaUB5E5C9ClUcMwb52GO0ay2Y9v1uBv1a4WhI3peKktAzYNk0EBMQlJtXPjRMrC9ylBPh+DsBHMu+KmDnfk7PIwyN4efC8k5kSRuPWoNdme1rz2+umU8FSmaWTHIajrbspf4GQbsntA5kuKEtDbfoNCU97o2KrRnUbeg3a8hwSjfh3u6MhlnGcg5K2Ij+zivEsWGCLKYUtE1ErqwfIzwWmJ6jnV66XCQGHf4Q1iIxqF7s2a1q24cgG2Z/iDXfqXrCIfy4P7b/Ztak3bdT9jfAdVZtdO5/r7I+O5hYhF86ayFlDWzZWP/ByiSb+q4CQbfVgK3BMmiAv2MqLHdhesmD/SmIcoOWUF6rFmRKZVFFpKpt5ATNTgUJ3JRowoXrrDruVXClUGRiCS6Zabd1rZ3VmTchaPJwtzQMdfIWISXj+Ig+C4UK0=",
		PullSecret: `{"auths":{"fake":{"auth":"Zm9vOmJhcg=="}}}`,
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
			},
		},
		Networking: &types.Networking{
			NetworkType: "OVNKubernetes",
			ClusterNetwork: []types.ClusterNetworkEntry{
				{
					CIDR:       *clusterNetwork,
					HostPrefix: 23,
				},
			},
			ServiceNetwork: []ipnet.IPNet{
				*serviceNetwork,
			},
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: *machineNetwork,
				},
			},
		},
		ControlPlane: &types.MachinePool{
			Name:     "master",
			Replicas: pointer.Int64Ptr(3),
		},
		Compute: []types.MachinePool{
			{
				Name:     "worker",
				Replicas: pointer.Int64Ptr(3),
			},
		},
	}

	// Set defaults
	defaults.SetInstallConfigDefaults(config)

	return config
}

func runCreateInstallConfig(ctx context.Context, req mcp.CallToolRequest) (string, error) {
	logrus.Info("MCP Server Creating install-config")
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

	// Set the install directory
	cluster.InstallDir = command.RootOpts.Dir

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      0,
				"total":         1,
				"progressToken": progressToken,
				"message":       "Creating install-config...",
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	// Create install-config
	installConfigRunner := runTargetCmd(ctx, targetassets.InstallConfig...)
	installConfigRunner(installConfigTarget.command, []string{})

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      1,
				"total":         1,
				"progressToken": progressToken,
				"message":       "Install-config created successfully",
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	return "Install-config created successfully", nil
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

	// Set the install directory
	cluster.InstallDir = command.RootOpts.Dir

	i := 0
	steps := 3

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      i,
				"total":         steps,
				"progressToken": progressToken,
				"message":       "Creating install-config...",
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if progressToken != nil {
					srv.SendNotificationToClient(
						ctx,
						"notifications/progress",
						map[string]any{
							"progress":      i,
							"total":         steps,
							"progressToken": progressToken,
							"message":       "Still working...",
						},
					)
				}
			case <-done:
				return
			}
		}
	}()

	// Create install-config first
	installConfigRunner := runTargetCmd(ctx, targetassets.InstallConfig...)
	installConfigRunner(installConfigTarget.command, []string{})

	i = 1
	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      i,
				"total":         steps,
				"progressToken": progressToken,
				"message":       "Creating cluster...",
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	// Now create the cluster
	runCommand := runTargetCmd(ctx, targetassets.Cluster...)
	runCommand(clusterTarget.command, []string{})

	i = 2
	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      i,
				"total":         steps,
				"progressToken": progressToken,
				"message":       "Finalizing cluster creation...",
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

	close(done)
	return "Cluster created successfully", nil
}

// createCustomInstallConfig creates an install-config based on parameters from the MCP client
func createCustomInstallConfig(params map[string]interface{}) (*types.InstallConfig, error) {
	clusterNetwork, _ := ipnet.ParseCIDR("10.128.0.0/14")
	serviceNetwork, _ := ipnet.ParseCIDR("172.30.0.0/16")
	machineNetwork, _ := ipnet.ParseCIDR("10.0.0.0/16")

	config := &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: params["cluster_name"].(string),
		},
		BaseDomain: params["base_domain"].(string),
		SSHKey:     params["ssh_key"].(string),
		PullSecret: params["pull_secret"].(string),
		Networking: &types.Networking{
			NetworkType: "OVNKubernetes",
			ClusterNetwork: []types.ClusterNetworkEntry{
				{
					CIDR:       *clusterNetwork,
					HostPrefix: 23,
				},
			},
			ServiceNetwork: []ipnet.IPNet{
				*serviceNetwork,
			},
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: *machineNetwork,
				},
			},
		},
		ControlPlane: &types.MachinePool{
			Name:     "master",
			Replicas: pointer.Int64Ptr(3),
		},
		Compute: []types.MachinePool{
			{
				Name:     "worker",
				Replicas: pointer.Int64Ptr(3),
			},
		},
	}

	// Set platform based on input
	platform := params["platform"].(string)
	region := params["region"].(string)

	switch platform {
	case "aws":
		config.Platform = types.Platform{
			AWS: &aws.Platform{
				Region: region,
			},
		}
	case "azure":
		config.Platform = types.Platform{
			Azure: &azure.Platform{
				Region: region,
			},
		}
	case "gcp":
		config.Platform = types.Platform{
			GCP: &gcp.Platform{
				Region: region,
			},
		}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	// Set defaults
	defaults.SetInstallConfigDefaults(config)

	return config, nil
}

func runCreateCustomInstallConfig(ctx context.Context, req mcp.CallToolRequest) (string, error) {
	logrus.Info("MCP Server Creating custom install-config")
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

	// Parse parameters from the request
	params, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid parameters format")
	}

	// Create custom install-config
	installConfig, err := createCustomInstallConfig(params)
	if err != nil {
		return "", fmt.Errorf("failed to create install-config: %w", err)
	}

	// Set the install directory
	cluster.InstallDir = command.RootOpts.Dir

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      0,
				"total":         1,
				"progressToken": progressToken,
				"message":       "Creating custom install-config...",
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	// Write the install-config to the target directory
	installConfigYAML, err := yaml.Marshal(installConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal install-config: %w", err)
	}

	installConfigPath := filepath.Join(command.RootOpts.Dir, "install-config.yaml")
	err = os.WriteFile(installConfigPath, installConfigYAML, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write install-config: %w", err)
	}

	if progressToken != nil {
		err := srv.SendNotificationToClient(
			ctx,
			"notifications/progress",
			map[string]any{
				"progress":      1,
				"total":         1,
				"progressToken": progressToken,
				"message":       "Custom install-config created successfully",
			},
		)
		if err != nil {
			logrus.Warn(err)
		}
	}

	return fmt.Sprintf("Custom install-config created successfully at %s", installConfigPath), nil
}
