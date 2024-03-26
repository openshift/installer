package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/explain"
)

func newExplainCmd(ctx context.Context) *cobra.Command {
	return explain.NewCmd()
}
