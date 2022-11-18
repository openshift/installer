package release

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

func NewCmd(f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Tools for managing the OpenShift release process",
		Long: templates.LongDesc(`
			This tool is used by OpenShift release to build images that can update a cluster.

			The subcommands allow you to see information about releases, perform administrative
			actions inspect the content of the release, and mirror release content across image
			registries.
			`),
	}
	cmd.AddCommand(NewInfo(f, streams))
	cmd.AddCommand(NewRelease(f, streams))
	cmd.AddCommand(NewExtract(f, streams))
	cmd.AddCommand(NewMirror(f, streams))
	return cmd
}
