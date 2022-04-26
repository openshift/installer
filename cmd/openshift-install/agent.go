package main

import (
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/agent"
)

func newAgentCmd() *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Commands for supporting cluster installation using agent installer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	agentCmd.AddCommand(agent.NewCreateCmd())
	agentCmd.AddCommand(agent.NewWaitForCmd())
	return agentCmd
}
