package network

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/constants"
	"github.com/openshift/rosa/pkg/reporter"
)

const (
	use     = "network"
	short   = "Network AWS cloudformation stack"
	long    = "Network AWS cloudformation stack using predefined yaml templates. "
	example = `  # Create a AWS cloudformation stack
  rosa create network <template-name> --param Param1=Value1 --param Param2=Value2 ` +
		"\n\n" + `  # ROSA quick start HCP VPC example with one availability zone` +
		"\n" + `  rosa create network rosa-quickstart-default-vpc --param Region=us-west-2` +
		` --param Name=quickstart-stack --param AvailabilityZoneCount=1` +
		` --param VpcCidr=10.0.0.0/16` +
		"\n\n" + `  # ROSA quick start HCP VPC example with two explicit availability zones` +
		"\n" + `  rosa create network rosa-quickstart-default-vpc --param Region=us-west-2` +
		` --param Name=quickstart-stack` +
		` --param AZ1=us-west-2b --param AZ2=us-west-2d --param VpcCidr=10.0.0.0/16` +
		"\n\n" + `  # To delete the AWS cloudformation stack` +
		"\n" + `  aws cloudformation delete-stack --stack-name <name> --region <region>` +
		"\n\n" + `# TEMPLATE_NAME:` +
		"\n" + `Specifies the name of the template to use. This should match the name of a directory ` +
		"\n" + `under the path specified by '--template-dir' or the 'OCM_TEMPLATE_DIR' environment variable.` +
		"\n" + `The directory should contain a YAML file defining the custom template structure.` +
		"\n\n" + `If no TEMPLATE_NAME is provided, or if no matching directory is found, the default ` +
		"\n" + `built-in template 'rosa-quickstart-default-vpc' will be used.`
	DefaultTemplateDir = "cmd/create/network/templates"
)

type NetworkUserOptions struct {
	Params      []string
	TemplateDir string
}

type NetworkOptions struct {
	reporter reporter.Logger
	args     *NetworkUserOptions
}

func NewNetworkUserOptions() *NetworkUserOptions {
	options := &NetworkUserOptions{}

	// Set template directory from environment variable or use default
	templateDir, isSet := os.LookupEnv(constants.OcmTemplateDir)
	if isSet {
		if templateDir == "\"\"" {
			templateDir = ""
		}
		options.TemplateDir = templateDir
	} else {
		options.TemplateDir = DefaultTemplateDir
	}

	return options
}

func (n *NetworkUserOptions) CleanTemplateDir() {
	// Clean up trailing '/' to work with filepath logic later
	if len(n.TemplateDir) > 0 && n.TemplateDir[len(n.TemplateDir)-1] == '/' {
		n.TemplateDir = n.TemplateDir[:len(n.TemplateDir)-1]
	}
}

func NewNetworkOptions() *NetworkOptions {
	return &NetworkOptions{
		reporter: reporter.CreateReporter(),
		args:     NewNetworkUserOptions(),
	}
}

func (m *NetworkOptions) Network() *NetworkUserOptions {
	return m.args
}

func BuildNetworkCommandWithOptions() (*cobra.Command, *NetworkUserOptions) {
	options := NewNetworkUserOptions()
	cmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Aliases: []string{"networks"},
		Example: example,
		Args:    cobra.MaximumNArgs(1),
		Hidden:  false,
	}
	var exportedTempDir string
	if options.TemplateDir != "" {
		exportedTempDir = options.TemplateDir
	}

	flags := cmd.Flags()
	flags.StringVar(
		&options.TemplateDir,
		"template-dir",
		"",
		"Use a specific template directory, overriding the OCM_TEMPLATE_DIR environment variable.",
	)
	if exportedTempDir != "" {
		options.TemplateDir = exportedTempDir
	}
	flags.StringArrayVar(
		&options.Params,
		"param",
		[]string{},
		"List of parameters",
	)

	return cmd, options
}
