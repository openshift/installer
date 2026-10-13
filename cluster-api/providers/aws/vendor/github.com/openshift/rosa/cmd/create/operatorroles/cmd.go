/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operatorroles

import (
	"fmt"
	"os"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/arguments"
	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/interactive"
	"github.com/openshift/rosa/pkg/interactive/confirm"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/roles"
	"github.com/openshift/rosa/pkg/rosa"
)

const (
	PrefixFlag             = "prefix"
	HostedCpFlag           = "hosted-cp"
	OidcConfigIdFlag       = "oidc-config-id"
	InstallerRoleArnFlag   = "role-arn"
	vpcEndpointRoleArnFlag = "vpc-endpoint-role-arn"
	hostedZoneRoleArnFlag  = "route53-role-arn"
)

var args struct {
	prefix              string
	hostedCp            bool
	installerRoleArn    string
	permissionsBoundary string
	forcePolicyCreation bool
	oidcConfigId        string
	sharedVpcRoleArn    string
	channelGroup        string
	vpcEndpointRoleArn  string
}

var Cmd = &cobra.Command{
	Use:     "operator-roles",
	Aliases: []string{"operatorroles"},
	Short:   "Create operator IAM roles for a cluster.",
	Long:    "Create cluster-specific operator IAM roles based on your cluster configuration.",
	Example: `  # Create default operator roles for cluster named "mycluster"
  rosa create operator-roles --cluster=mycluster

  # Create operator roles with a specific permissions boundary
  rosa create operator-roles -c mycluster --permissions-boundary arn:aws:iam::123456789012:policy/perm-boundary`,
	Run:  run,
	Args: cobra.MaximumNArgs(3),
}

func init() {
	flags := Cmd.Flags()

	ocm.AddOptionalClusterFlag(Cmd)

	flags.StringVar(
		&args.prefix,
		PrefixFlag,
		"",
		"User-defined prefix for generated AWS operator policies. Not to be used alongside --cluster flag.",
	)

	flags.StringVar(
		&args.oidcConfigId,
		OidcConfigIdFlag,
		"",
		"Registered OIDC configuration ID to add its issuer URL as the trusted relationship to the operator roles. "+
			"Not to be used alongside --cluster flag.",
	)

	// normalizing installer role argument to support deprecated flag
	flags.SetNormalizeFunc(arguments.NormalizeFlags)
	flags.StringVar(
		&args.installerRoleArn,
		InstallerRoleArnFlag,
		"",
		"Installer role ARN supplied to retrieve operator policy prefix and path. Not to be used alongside --cluster flag.",
	)

	flags.BoolVar(
		&args.hostedCp,
		HostedCpFlag,
		false,
		"Indicates whether to create the hosted control planes operator roles when using --prefix option.",
	)

	flags.StringVar(
		&args.permissionsBoundary,
		"permissions-boundary",
		"",
		"The ARN of the policy that is used to set the permissions boundary for the operator roles.",
	)

	flags.BoolVarP(
		&args.forcePolicyCreation,
		"force-policy-creation",
		"f",
		false,
		"Forces creation of policies skipping compatibility check",
	)

	flags.StringVar(
		&args.sharedVpcRoleArn,
		"shared-vpc-role-arn",
		"",
		"AWS IAM role ARN with a policy attached, granting permissions necessary to create and manage Route 53 DNS records "+
			"in private Route 53 hosted zone associated with intended shared VPC.",
	)

	flags.StringVar(
		&args.channelGroup,
		"channel-group",
		ocm.DefaultChannelGroup,
		"Channel group is the name of the channel where this image belongs, for example \"stable\" or \"fast\".",
	)
	flags.MarkHidden("channel-group")

	flags.StringVar(
		&args.vpcEndpointRoleArn,
		vpcEndpointRoleArnFlag,
		"",
		"AWS IAM Role ARN with policy attached, associated with the shared VPC."+
			" Grants permissions necessary to communicate with and handle a Hosted Control Plane cross-account VPC.",
	)

	flags.StringVar(
		&args.sharedVpcRoleArn,
		hostedZoneRoleArnFlag,
		"",
		"AWS IAM Role Arn with policy attached, associated with shared VPC."+
			" Grants permission necessary to handle route53 operations associated with a cross-account VPC. "+
			"This flag deprecates '--shared-vpc-role-arn'.",
	)

	flags.MarkDeprecated("shared-vpc-role-arn", fmt.Sprintf("'--shared-vpc-role-arn' will be replaced with "+
		"'--%s' in future versions of ROSA.", hostedZoneRoleArnFlag))

	interactive.AddModeFlag(Cmd)
	confirm.AddFlag(flags)
	interactive.AddFlag(flags)
}

func run(cmd *cobra.Command, argv []string) {
	r := rosa.NewRuntime().WithAWS().WithOCM()
	defer r.Cleanup()

	rosa.HostedClusterOnlyFlag(r, cmd, vpcEndpointRoleArnFlag)

	// Allow the command to be called programmatically
	isProgmaticallyCalled := false
	if len(argv) == 3 && !cmd.Flag("cluster").Changed {
		ocm.SetClusterKey(argv[0])
		interactive.SetModeKey(argv[1])
		args.permissionsBoundary = argv[2]

		// if mode is empty skip interactive is true
		if argv[1] != "" {
			isProgmaticallyCalled = true
		}
	}

	var isHcpSharedVpc bool
	var err error
	if !args.hostedCp {
		rosa.HostedClusterOnlyFlag(r, cmd, vpcEndpointRoleArnFlag)
	} else {
		isHcpSharedVpc, err = roles.ValidateSharedVpcInputs(args.vpcEndpointRoleArn, args.sharedVpcRoleArn,
			vpcEndpointRoleArnFlag, hostedZoneRoleArnFlag)
		if err != nil {
			r.Reporter.Errorf("%s", err)
			os.Exit(1)
		}
	}

	if args.vpcEndpointRoleArn != "" {
		err = aws.ARNValidator(args.vpcEndpointRoleArn)
		if err != nil {
			r.Reporter.Errorf("Expected a valid policy ARN for %s: %s", vpcEndpointRoleArnFlag, err)
			os.Exit(1)
		}
	}
	if args.sharedVpcRoleArn != "" {
		err = aws.ARNValidator(args.sharedVpcRoleArn)
		if err != nil {
			r.Reporter.Errorf("Expected a valid policy ARN for %s: %s", hostedZoneRoleArnFlag, err)
			os.Exit(1)
		}
	}

	env, err := ocm.GetEnv()
	if err != nil {
		r.Reporter.Errorf("Failed to determine OCM environment: %v", err)
		os.Exit(1)
	}

	mode, err := interactive.GetMode()
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}

	// Determine if interactive mode is needed
	if !interactive.Enabled() && !cmd.Flags().Changed("mode") && !isProgmaticallyCalled {
		interactive.Enable()
	}

	if !cmd.Flag("cluster").Changed && !cmd.Flag(PrefixFlag).Changed && !isProgmaticallyCalled {
		r.Reporter.Errorf("Either a cluster key for STS cluster or an operator roles prefix must be specified.")
		os.Exit(1)
	}

	if cmd.Flag("cluster").Changed && cmd.Flag(PrefixFlag).Changed {
		r.Reporter.Errorf("A cluster key for STS cluster and an operator roles prefix " +
			"cannot be specified alongside each other.")
		os.Exit(1)
	}

	if cmd.Flag("cluster").Changed && cmd.Flag(OidcConfigIdFlag).Changed {
		r.Reporter.Errorf("A cluster key for STS cluster and an OIDC configuration ID " +
			"cannot be specified alongside each other.")
		os.Exit(1)
	}

	if !args.hostedCp && args.installerRoleArn != "" {
		managedPolicies, err := r.AWSClient.HasManagedPolicies(args.installerRoleArn)
		if err != nil {
			r.Reporter.Errorf("Failed to determine if cluster has managed policies: %v", err)
			os.Exit(1)
		}
		if managedPolicies {
			r.Reporter.Errorf("The managed policies are not supported for classic operator-roles.")
			os.Exit(1)
		}
	}

	var cluster *cmv1.Cluster
	if args.prefix == "" {
		cluster = r.FetchCluster()
	}

	if args.forcePolicyCreation && mode != interactive.ModeAuto {
		r.Reporter.Warnf("Forcing creation of policies only works in auto mode")
		os.Exit(1)
	}

	if interactive.Enabled() && !isProgmaticallyCalled {
		mode, err = interactive.GetOptionMode(cmd, mode, "Role creation mode")
		if err != nil {
			r.Reporter.Errorf("Expected a valid role creation mode: %s", err)
			os.Exit(1)
		}
	}

	if cluster == nil && interactive.Enabled() && !isProgmaticallyCalled {
		handleOperatorRolesPrefixOptions(r, cmd)
	}

	permissionsBoundary := args.permissionsBoundary
	if interactive.Enabled() && !isProgmaticallyCalled {
		permissionsBoundary, err = interactive.GetString(interactive.Input{
			Question: "Permissions boundary ARN",
			Help:     cmd.Flags().Lookup("permissions-boundary").Usage,
			Default:  permissionsBoundary,
			Validators: []interactive.Validator{
				aws.ARNValidator,
			},
		})
		if err != nil {
			r.Reporter.Errorf("Expected a valid policy ARN for permissions boundary: %s", err)
			os.Exit(1)
		}
	}

	if interactive.Enabled() && args.hostedCp {
		defaultValue := args.sharedVpcRoleArn != "" && args.vpcEndpointRoleArn != ""
		isHcpSharedVpc, err = interactive.GetBool(interactive.Input{
			Question: "Use operator roles for Hosted CP shared VPC?",
			Help: "Whether or not to set route53/VPC endpoint role ARNs to be used for Hosted CP shared VPC " +
				"(cross-account VPC)",
			Default:  defaultValue,
			Required: false,
		})
		if err != nil {
			r.Reporter.Errorf("Expected a valid value: %s", err)
			os.Exit(1)
		}

		if !isHcpSharedVpc {
			args.sharedVpcRoleArn = ""
			args.vpcEndpointRoleArn = ""
		}
	}

	if interactive.Enabled() && isHcpSharedVpc && !r.Creator.IsGovcloud {
		args.vpcEndpointRoleArn, err = interactive.GetString(interactive.Input{
			Question: "Set VPC endpoint role ARN",
			Help:     cmd.Flags().Lookup(vpcEndpointRoleArnFlag).Usage,
			Default:  args.vpcEndpointRoleArn,
			Required: isHcpSharedVpc,
			Validators: []interactive.Validator{
				aws.ARNValidator,
			},
		})
		if err != nil {
			r.Reporter.Errorf("Expected a valid value: %s", err)
			os.Exit(1)
		}
	}
	if interactive.Enabled() && isHcpSharedVpc && !r.Creator.IsGovcloud {
		args.sharedVpcRoleArn, err = interactive.GetString(interactive.Input{
			Question: "Set route53 role ARN",
			Help:     cmd.Flags().Lookup(hostedZoneRoleArnFlag).Usage,
			Default:  args.sharedVpcRoleArn,
			Required: isHcpSharedVpc,
			Validators: []interactive.Validator{
				aws.ARNValidator,
			},
		})
		if err != nil {
			r.Reporter.Errorf("Expected a valid value: %s", err)
			os.Exit(1)
		}
	}

	if permissionsBoundary != "" {
		err = aws.ARNValidator(permissionsBoundary)
		if err != nil {
			r.Reporter.Errorf("Expected a valid policy ARN for permissions boundary: %s", err)
			os.Exit(1)
		}
	}

	policies, err := r.OCMClient.GetPolicies("OperatorRole")
	if err != nil {
		r.Reporter.Errorf("Expected a valid role creation mode: %s", err)
		os.Exit(1)
	}

	if args.prefix != "" {
		if args.oidcConfigId == "" {
			r.Reporter.Errorf("%s is mandatory for %s param flow.", OidcConfigIdFlag, PrefixFlag)
			os.Exit(1)
		}

		if args.installerRoleArn == "" {
			r.Reporter.Errorf("%s is mandatory for %s param flow.", InstallerRoleArnFlag, PrefixFlag)
			os.Exit(1)
		}
		channelGroup := args.channelGroup
		latestPolicyVersion, err := r.OCMClient.GetLatestVersion(channelGroup)
		if err != nil {
			r.Reporter.Errorf("Error getting latest version: %s", err)
			os.Exit(1)
		}
		err = HandleOperatorRoleCreationByPrefix(r, env, permissionsBoundary,
			mode, policies, latestPolicyVersion, isHcpSharedVpc)
		if err != nil {
			r.Reporter.Errorf("Error creating operator roles: %s", err)
			os.Exit(1)
		}
		return
	}
	latestPolicyVersion, err := r.OCMClient.GetLatestVersion(cluster.Version().ChannelGroup())
	if err != nil {
		r.Reporter.Errorf("Error getting latest version: %s", err)
		os.Exit(1)
	}
	err = handleOperatorRoleCreationByClusterKey(r, env, permissionsBoundary,
		mode, policies, latestPolicyVersion, isHcpSharedVpc)
	if err != nil {
		r.Reporter.Errorf("Error creating operator roles: %s", err)
		os.Exit(1)
	}
}

func convertV1OperatorIAMRoleIntoOcmOperatorIamRole(
	operatorIAMRoleList []*cmv1.OperatorIAMRole) ([]ocm.OperatorIAMRole, error) {
	operatorRolesList := []ocm.OperatorIAMRole{}
	for _, operatorIAMRole := range operatorIAMRoleList {
		newRole, err := ocm.NewOperatorIamRoleFromCmv1(operatorIAMRole)
		if err != nil {
			return operatorRolesList, err
		}
		operatorRolesList = append(operatorRolesList, *newRole)
	}
	return operatorRolesList, nil
}
