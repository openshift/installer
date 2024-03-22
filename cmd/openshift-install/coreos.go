package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/coreoscli"
)

func newCoreOSCmd(ctx context.Context) *cobra.Command {
	return coreoscli.NewCmd()
}
