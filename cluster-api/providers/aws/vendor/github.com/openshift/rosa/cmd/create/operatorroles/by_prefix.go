package operatorroles

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	awserr "github.com/openshift-online/ocm-common/pkg/aws/errors"
	common "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/aws"
	awscb "github.com/openshift/rosa/pkg/aws/commandbuilder"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/interactive"
	interactiveOidc "github.com/openshift/rosa/pkg/interactive/oidc"
	interactiveRoles "github.com/openshift/rosa/pkg/interactive/roles"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/output"
	"github.com/openshift/rosa/pkg/roles"
	"github.com/openshift/rosa/pkg/rosa"
)

func handleOperatorRolesPrefixOptions(r *rosa.Runtime, cmd *cobra.Command) {
	operatorRolesPrefix := args.prefix
	operatorRolesPrefix, err := interactive.GetString(interactive.Input{
		Question: "Operator roles prefix",
		Help:     cmd.Flags().Lookup(PrefixFlag).Usage,
		Required: true,
		Default:  operatorRolesPrefix,
		Validators: []interactive.Validator{
			interactive.RegExp(aws.RoleNameRE.String()),
			interactive.MaxLength(32),
		},
	})
	if err != nil {
		r.Reporter.Errorf("Expected a prefix for the operator IAM roles: %s", err)
		os.Exit(1)
	}
	args.prefix = operatorRolesPrefix

	if args.oidcConfigId == "" {
		args.oidcConfigId = interactiveOidc.GetOidcConfigID(r, cmd)
	}

	isHostedCP := args.hostedCp
	isHostedCP, err = interactive.GetBool(interactive.Input{
		Question: "Create hosted control plane operator roles",
		Help:     cmd.Flags().Lookup("hosted-cp").Usage,
		Default:  isHostedCP,
		Required: false,
	})
	if err != nil {
		_ = r.Reporter.Errorf("Expected a valid --hosted-cp value: %s", err)
		os.Exit(1)
	}

	args.hostedCp = isHostedCP
	if args.hostedCp {
		args.installerRoleArn = interactiveRoles.GetInstallerRoleArn(r, cmd, args.installerRoleArn, "",
			r.AWSClient.FindRoleARNsHostedCp)
	} else {
		args.installerRoleArn = interactiveRoles.GetInstallerRoleArn(r, cmd, args.installerRoleArn, "",
			r.AWSClient.FindRoleARNsClassic)
	}
}

func CreateOperatorRoles(r *rosa.Runtime, env string, permissionsBoundary string, mode string, policies map[string]*cmv1.AWSSTSPolicy, defaultPolicyVersion string, isSharedVpc bool,
	prefix string, hostedCp bool, installerRoleArn string, forcePolicyCreation bool, oidcConfigId string, sharedVpcRoleArn string, channelGroup string, vpcEndpointRoleArn string) error {
	args.prefix = prefix
	args.hostedCp = hostedCp
	args.installerRoleArn = installerRoleArn
	args.permissionsBoundary = permissionsBoundary
	args.forcePolicyCreation = forcePolicyCreation
	args.oidcConfigId = oidcConfigId
	args.sharedVpcRoleArn = sharedVpcRoleArn
	args.channelGroup = channelGroup
	args.vpcEndpointRoleArn = vpcEndpointRoleArn

	return HandleOperatorRoleCreationByPrefix(r, env, permissionsBoundary, mode, policies, defaultPolicyVersion, isSharedVpc)
}

func HandleOperatorRoleCreationByPrefix(r *rosa.Runtime, env string,
	permissionsBoundary string, mode string,
	policies map[string]*cmv1.AWSSTSPolicy,
	defaultPolicyVersion string, isSharedVpc bool) error {
	oidcConfig, err := r.OCMClient.GetOidcConfig(args.oidcConfigId)
	if err != nil {
		r.Reporter.Errorf("There was a problem retrieving OIDC Config '%s': %v", args.oidcConfigId, err)
		os.Exit(1)
	}
	includeHostedCpSet := args.hostedCp
	operatorRolesPrefix := args.prefix
	oidcEndpointUrl := oidcConfig.IssuerUrl()
	installerRoleArn := args.installerRoleArn
	sharedVpcRoleArn := args.sharedVpcRoleArn
	sharedVpcEndpointRoleArn := args.vpcEndpointRoleArn

	validateArgumentsOperatorRolesCreationByPrefix(r, operatorRolesPrefix, oidcEndpointUrl, installerRoleArn)

	installerRoleName, err := aws.GetResourceIdFromARN(installerRoleArn)
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}
	path, err := aws.GetPathFromARN(installerRoleArn)
	if err != nil {
		r.Reporter.Errorf("Expected a valid path for '%s': %v", installerRoleArn, err)
		os.Exit(1)
	}
	if path != "" && !output.HasFlag() && r.Reporter.IsTerminal() {
		r.Reporter.Infof("ARN path '%s' detected in installer role '%s'. "+
			"This ARN path will be used for subsequent created operator roles and policies.",
			path, installerRoleArn)
	}

	hasStandardNamedInstallerRole, installerRolePrefix := aws.IsStandardNamedAccountRole(installerRoleName,
		aws.AccountRoles[aws.InstallerAccountRole].Name)
	if !hasStandardNamedInstallerRole {
		r.Reporter.Infof("Can only use installer roles created through ROSA CLI for this flow.")
		os.Exit(1)
	}
	operatorRolePolicyPrefix := installerRolePrefix
	credRequests, err := r.OCMClient.GetCredRequests(includeHostedCpSet)
	if err != nil {
		r.Reporter.Errorf("Error getting operator credential request from OCM %s", err)
		os.Exit(1)
	}
	managedPolicies, err := r.AWSClient.HasManagedPolicies(installerRoleArn)
	if err != nil {
		r.Reporter.Errorf("Failed to determine if cluster has managed policies: %v", err)
		os.Exit(1)
	}
	awsCreator, err := r.AWSClient.GetCreator()
	if err != nil {
		r.Reporter.Errorf("Unable to get IAM credentials: %v", err)
		os.Exit(1)
	}

	operatorIAMRoleList, err := convertCredRequestsOperatorRolesIntoV1OperatorIAMRole(credRequests,
		args.prefix, awsCreator, path)
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}

	var hostedCPPolicies bool
	if args.hostedCp {
		hostedCPPolicies, err = r.AWSClient.HasHostedCPPolicies(args.installerRoleArn)
		if err != nil {
			r.Reporter.Errorf("Failed to determine if the Installer role ARN has hosted CP policies: %v", err)
			os.Exit(1)
		}

		if !hostedCPPolicies {
			r.Reporter.Errorf(
				"Failed to create the operator role since the Installer role ARN '%v' does not have managed policies",
				args.installerRoleArn)
			os.Exit(1)
		}
	}

	operatorRolesList, err := convertV1OperatorIAMRoleIntoOcmOperatorIamRole(operatorIAMRoleList)
	if err != nil {
		r.Reporter.Errorf("%v", err)
		os.Exit(1)
	}
	err = ocm.ValidateOperatorRolesMatchOidcProvider(r.Reporter, r.AWSClient,
		operatorRolesList, oidcConfig.IssuerUrl(), "4.0", path, managedPolicies, true)
	if err != nil && !awserr.IsNoSuchEntityException(err) {
		r.Reporter.Errorf("%v", err)
		os.Exit(1)
	}

	switch mode {
	case interactive.ModeAuto:
		if !output.HasFlag() || r.Reporter.IsTerminal() {
			r.Reporter.Infof("Creating roles using '%s'", r.Creator.ARN)
		}
		err = createRolesByPrefix(r, operatorRolePolicyPrefix, permissionsBoundary,
			defaultPolicyVersion, policies,
			credRequests, managedPolicies,
			path, operatorIAMRoleList,
			oidcEndpointUrl, hostedCPPolicies, sharedVpcRoleArn, sharedVpcEndpointRoleArn,
			isSharedVpc)
		if err != nil {
			r.Reporter.Errorf("There was an error creating the operator roles: %s", err)
			isThrottle := "false"
			if strings.Contains(err.Error(), "Throttling") {
				isThrottle = helper.True
			}
			r.OCMClient.LogEvent("ROSACreateOperatorRolesModeAuto", map[string]string{
				ocm.OperatorRolesPrefix: operatorRolesPrefix,
				ocm.Response:            ocm.Failure,
				ocm.IsThrottle:          isThrottle,
			})
			os.Exit(1)
		}
		if r.Reporter.IsTerminal() {
			hostedCpOutputParam := ""
			stsParam := " --sts"
			if args.hostedCp {
				hostedCpOutputParam = fmt.Sprintf(" --%s", HostedCpFlag)
				// HCP clusters are STS by default, so we don't need the --sts flag
				stsParam = ""
			}
			installerRoleArnParam := ""
			if args.installerRoleArn != "" && path != "" {
				installerRoleArnParam = fmt.Sprintf(" --role-arn %s", args.installerRoleArn)
			}
			r.Reporter.Infof(fmt.Sprintf("To create a cluster with these roles, run the following command:\n"+
				"\trosa create cluster%s --oidc-config-id %s --operator-roles-prefix %s%s%s",
				stsParam, args.oidcConfigId, args.prefix, hostedCpOutputParam, installerRoleArnParam))
		}
		r.OCMClient.LogEvent("ROSACreateOperatorRolesModeAuto", map[string]string{
			ocm.OperatorRolesPrefix: operatorRolesPrefix,
			ocm.Response:            ocm.Success,
		})
	case interactive.ModeManual:
		commands, err := buildCommandsFromPrefix(r, env,
			operatorRolePolicyPrefix, permissionsBoundary,
			defaultPolicyVersion, policies,
			credRequests, managedPolicies,
			path, operatorIAMRoleList,
			oidcEndpointUrl, hostedCPPolicies, sharedVpcRoleArn, sharedVpcEndpointRoleArn)
		if err != nil {
			r.Reporter.Errorf("There was an error building the list of resources: %s", err)
			os.Exit(1)
			r.OCMClient.LogEvent("ROSACreateOperatorRolesModeManual", map[string]string{
				ocm.OperatorRolesPrefix: operatorRolesPrefix,
				ocm.Response:            ocm.Failure,
			})
		}
		if r.Reporter.IsTerminal() {
			r.Reporter.Infof("All policy files saved to the current directory")
			r.Reporter.Infof("Run the following commands to create the operator roles:\n")
		}
		r.OCMClient.LogEvent("ROSACreateOperatorRolesModeManual", map[string]string{
			ocm.OperatorRolesPrefix: operatorRolesPrefix,
		})
		fmt.Println(commands)
	default:
		r.Reporter.Errorf("Invalid mode. Allowed values are %s", interactive.Modes)
		os.Exit(1)
	}
	return nil
}

func convertCredRequestsOperatorRolesIntoV1OperatorIAMRole(credRequests map[string]*cmv1.STSOperator,
	operatorRolesPrefix string, awsCreator *aws.Creator, path string) ([]*cmv1.OperatorIAMRole, error) {
	operatorIAMRoleList := []*cmv1.OperatorIAMRole{}
	for _, operator := range credRequests {
		operatorIamRole, err := cmv1.NewOperatorIAMRole().
			Name(operator.Name()).
			Namespace(operator.Namespace()).
			RoleARN(aws.ComputeOperatorRoleArn(operatorRolesPrefix, operator,
				awsCreator, path)).
			Build()
		if err != nil {
			return operatorIAMRoleList, err
		}
		operatorIAMRoleList = append(operatorIAMRoleList, operatorIamRole)
	}
	return operatorIAMRoleList, nil
}

func validateArgumentsOperatorRolesCreationByPrefix(r *rosa.Runtime, operatorRolesPrefix string,
	oidcEndpointUrl string, installerRoleArn string) {
	if len(operatorRolesPrefix) == 0 {
		r.Reporter.Errorf("Expected a prefix for the operator IAM roles")
		os.Exit(1)
	}
	if len(operatorRolesPrefix) > 32 {
		r.Reporter.Errorf("Expected a prefix with no more than 32 characters")
		os.Exit(1)
	}
	if !aws.RoleNameRE.MatchString(operatorRolesPrefix) {
		r.Reporter.Errorf("Expected valid operator roles prefix matching %s", aws.RoleNameRE.String())
		os.Exit(1)
	}
	parsedURI, err := url.ParseRequestURI(oidcEndpointUrl)
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}
	if parsedURI.Scheme != helper.ProtocolHttps {
		r.Reporter.Errorf("Expected OIDC endpoint URL '%s' to use an https:// scheme", oidcEndpointUrl)
		os.Exit(1)
	}
	err = aws.ARNValidator(installerRoleArn)
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}
}

func createRolesByPrefix(r *rosa.Runtime, prefix string, permissionsBoundary string, defaultPolicyVersion string,
	policies map[string]*cmv1.AWSSTSPolicy, credRequests map[string]*cmv1.STSOperator,
	managedPolicies bool, path string,
	operatorIAMRoleList []*cmv1.OperatorIAMRole,
	oidcEndpointUrl string, hostedCPPolicies bool, sharedVpcRoleArn string, sharedVpcEndpointRoleArn string,
	isHcpSharedVpc bool) error {

	isSharedVpc := sharedVpcRoleArn != ""

	for credrequest, operator := range credRequests {
		roleArn := aws.FindOperatorRoleBySTSOperator(operatorIAMRoleList, operator)
		roleName, err := aws.GetResourceIdFromARN(roleArn)
		if err != nil {
			return err
		}
		if roleName == "" {
			return fmt.Errorf("Failed to find operator IAM role")
		}

		var policyArn string
		var policyArns []string
		filename := aws.GetOperatorPolicyKey(credrequest, hostedCPPolicies, isSharedVpc)
		if managedPolicies {
			policyArn, err = aws.GetManagedPolicyARN(policies, filename)
			if err != nil {
				return err
			}
			if isHcpSharedVpc {
				if credrequest == aws.IngressOperatorCloudCredentialsRoleType {
					sharedVpcPolicyArn, err := getHcpSharedVpcPolicy(r, sharedVpcRoleArn, defaultPolicyVersion)
					if err != nil {
						return err
					}
					policyArns = append(policyArns, sharedVpcPolicyArn)
				} else if credrequest == aws.ControlPlaneCloudCredentialsRoleType {
					for _, arn := range []string{sharedVpcEndpointRoleArn, sharedVpcRoleArn} {
						sharedVpcPolicyArn, err := getHcpSharedVpcPolicy(r, arn, defaultPolicyVersion)
						if err != nil {
							return err
						}
						policyArns = append(policyArns, sharedVpcPolicyArn)
					}
				}
			}
		} else {
			policyArn = aws.GetOperatorPolicyARN(r.Creator.Partition, r.Creator.AccountID, prefix, operator.Namespace(),
				operator.Name(), path)
			policyDetails := aws.GetPolicyDetails(policies, filename)

			if isSharedVpc && credrequest == aws.IngressOperatorCloudCredentialsRoleType {
				err = validateIngressOperatorPolicyOverride(r, policyArn, sharedVpcRoleArn, prefix)
				if err != nil {
					return err
				}

				policyDetails = aws.InterpolatePolicyDocument(r.Creator.Partition, policyDetails, map[string]string{
					"shared_vpc_role_arn": sharedVpcRoleArn,
				})
			}

			operatorPolicyTags := map[string]string{
				common.OpenShiftVersion: defaultPolicyVersion,
				tags.RolePrefix:         prefix,
				tags.RedHatManaged:      helper.True,
				tags.OperatorNamespace:  operator.Namespace(),
				tags.OperatorName:       operator.Name(),
			}

			if args.forcePolicyCreation || (isSharedVpc && credrequest == aws.IngressOperatorCloudCredentialsRoleType) {
				_, err := r.AWSClient.ForceEnsurePolicy(policyArn, policyDetails,
					defaultPolicyVersion, operatorPolicyTags, path)
				if err != nil {
					return err
				}
			} else {
				_, err := r.AWSClient.EnsurePolicy(policyArn, policyDetails,
					defaultPolicyVersion, operatorPolicyTags, path)
				if err != nil {
					return err
				}
			}
		}
		policyArns = append(policyArns, policyArn)

		policyDetails := aws.GetPolicyDetails(policies, "operator_iam_role_policy")
		policy, err := aws.GenerateOperatorRolePolicyDocByOidcEndpointUrl(r.Creator.Partition, oidcEndpointUrl,
			r.Creator.AccountID, operator, policyDetails)
		if err != nil {
			return err
		}

		r.Reporter.Debugf("Creating role '%s'", roleName)
		tagsList := map[string]string{
			tags.OperatorNamespace: operator.Namespace(),
			tags.OperatorName:      operator.Name(),
			tags.RedHatManaged:     helper.True,
		}
		if managedPolicies {
			tagsList[common.ManagedPolicies] = helper.True
		}
		if hostedCPPolicies {
			tagsList[tags.HypershiftPolicies] = helper.True
		}

		roleARN, err := r.AWSClient.EnsureRole(r.Reporter, roleName, policy, permissionsBoundary, defaultPolicyVersion,
			tagsList, path, managedPolicies)
		if err != nil {
			return err
		}
		if !output.HasFlag() || r.Reporter.IsTerminal() {
			r.Reporter.Infof("Created role '%s' with ARN '%s'", roleName, roleARN)
		}

		for _, arn := range policyArns {
			r.Reporter.Debugf("Attaching permission policy '%s' to role '%s'", arn, roleName)
			err = r.AWSClient.AttachRolePolicy(r.Reporter, roleName, arn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func buildCommandsFromPrefix(r *rosa.Runtime, env string,
	prefix string, permissionsBoundary string, defaultPolicyVersion string,
	policies map[string]*cmv1.AWSSTSPolicy, credRequests map[string]*cmv1.STSOperator,
	managedPolicies bool, path string,
	operatorIAMRoleList []*cmv1.OperatorIAMRole,
	oidcEndpointUrl string, hostedCPPolicies bool, sharedVpcRoleArn string, vpcEndpointRoleArn string) (string, error) {
	if !managedPolicies {
		err := aws.GenerateOperatorRolePolicyFiles(r.Reporter, policies, credRequests, sharedVpcRoleArn, r.Creator.Partition)
		if err != nil {
			r.Reporter.Errorf("There was an error generating the policy files: %s", err)
			os.Exit(1)
		}
	}

	isSharedVpc := sharedVpcRoleArn != ""
	var policyDetails = make(map[string]roles.ManualSharedVpcPolicyDetails)

	commands := []string{}

	for credrequest, operator := range credRequests {
		roleArn := aws.FindOperatorRoleBySTSOperator(operatorIAMRoleList, operator)
		roleName, err := aws.GetResourceIdFromARN(roleArn)
		if err != nil {
			return "", err
		}

		var policyARN string
		if managedPolicies {
			policyARN, err = aws.GetManagedPolicyARN(policies, aws.GetOperatorPolicyKey(
				credrequest, hostedCPPolicies, false))
			if err != nil {
				return "", err
			}
		} else {
			policyARN = computePolicyARN(*r.Creator, prefix, operator.Namespace(), operator.Name(), path)
			name := aws.GetOperatorPolicyName(prefix, operator.Namespace(), operator.Name())
			iamTags := map[string]string{
				common.OpenShiftVersion: defaultPolicyVersion,
				tags.RolePrefix:         prefix,
				tags.OperatorNamespace:  operator.Namespace(),
				tags.OperatorName:       operator.Name(),
				tags.RedHatManaged:      helper.True,
			}
			operatorPolicyKey := aws.GetOperatorPolicyKey(credrequest, hostedCPPolicies, isSharedVpc)
			fileName := fmt.Sprintf("file://%s.json", operatorPolicyKey)
			_, err = r.AWSClient.IsPolicyExists(policyARN)
			if err != nil {
				createPolicy := awscb.NewIAMCommandBuilder().
					SetCommand(awscb.CreatePolicy).
					AddParam(awscb.PolicyName, name).
					AddParam(awscb.PolicyDocument, fileName).
					AddTags(iamTags).
					AddParam(awscb.Path, path).
					Build()
				commands = append(commands, createPolicy)
			} else if isSharedVpc && credrequest == aws.IngressOperatorCloudCredentialsRoleType {
				err := validateIngressOperatorPolicyOverride(r, policyARN, sharedVpcRoleArn, prefix)
				if err != nil {
					return "", err
				}

				createPolicyVersion := awscb.NewIAMCommandBuilder().
					SetCommand(awscb.CreatePolicyVersion).
					AddParam(awscb.PolicyArn, policyARN).
					AddParam(awscb.PolicyDocument, fileName).
					AddParamNoValue(awscb.SetAsDefault).
					Build()
				commands = append(commands, createPolicyVersion)
			}
		}

		policyDetail := aws.GetPolicyDetails(policies, "operator_iam_role_policy")
		policy, err := aws.GenerateOperatorRolePolicyDocByOidcEndpointUrl(r.Creator.Partition, oidcEndpointUrl,
			r.Creator.AccountID, operator, policyDetail)
		if err != nil {
			return "", err
		}

		filename := fmt.Sprintf("operator_%s_policy", credrequest)
		filename = aws.GetFormattedFileName(filename)
		r.Reporter.Debugf("Saving '%s' to the current directory", filename)
		err = helper.SaveDocument(policy, filename)
		if err != nil {
			return "", err
		}
		iamTags := map[string]string{
			tags.OperatorNamespace: operator.Namespace(),
			tags.OperatorName:      operator.Name(),
			tags.RedHatManaged:     helper.True,
		}
		if managedPolicies {
			iamTags[common.ManagedPolicies] = helper.True
		}
		if hostedCPPolicies {
			iamTags[tags.HypershiftPolicies] = helper.True
		}
		createRole := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.CreateRole).
			AddParam(awscb.RoleName, roleName).
			AddParam(awscb.AssumeRolePolicyDocument, fmt.Sprintf("file://%s", filename)).
			AddParam(awscb.PermissionsBoundary, permissionsBoundary).
			AddTags(iamTags).
			AddParam(awscb.Path, path).
			Build()

		attachRolePolicy := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.AttachRolePolicy).
			AddParam(awscb.RoleName, roleName).
			AddParam(awscb.PolicyArn, policyARN).
			Build()

		var attachSharedVpcRolePolicy string
		var policyCommands []string

		if isSharedVpc { // HCP Shared VPC policy attachment

			// Precreate HCP shared VPC policies for less memory usage + time to execute
			// Shared VPC role arn (route53)
			if _, ok := policyDetails[aws.IngressOperatorCloudCredentialsRoleType]; !ok {
				exists, createPolicyCommand, policyName, err := roles.GetHcpSharedVpcPolicyDetails(r, sharedVpcRoleArn)
				if err != nil {
					return "", err
				}

				sharedVpcRolePath, err := aws.GetPathFromARN(sharedVpcRoleArn)
				if err != nil {
					return "", err
				}

				policyDetails[aws.IngressOperatorCloudCredentialsRoleType] = roles.ManualSharedVpcPolicyDetails{
					Command:       createPolicyCommand,
					Name:          policyName,
					AlreadyExists: exists,
					Path:          sharedVpcRolePath,
				}
			}
			// VPC endpoint role arn
			if _, ok := policyDetails[aws.ControlPlaneCloudCredentialsRoleType]; !ok {

				exists, createPolicyCommand, policyName, err := roles.GetHcpSharedVpcPolicyDetails(r, vpcEndpointRoleArn)
				if err != nil {
					return "", err
				}

				vpcEndpointRolePath, err := aws.GetPathFromARN(vpcEndpointRoleArn)
				if err != nil {
					return "", err
				}

				policyDetails[aws.ControlPlaneCloudCredentialsRoleType] = roles.ManualSharedVpcPolicyDetails{
					Command:       createPolicyCommand,
					Name:          policyName,
					AlreadyExists: exists,
					Path:          vpcEndpointRolePath,
				}
			}

			var policies []string

			// Attach HCP shared VPC policies
			switch credrequest {
			case aws.IngressOperatorCloudCredentialsRoleType:
				if details, ok := policyDetails[credrequest]; ok {
					policies = append(policies, policyDetails[credrequest].Name)
					if !policyDetails[credrequest].AlreadyExists { // Skip creation if already exists
						policyCommands = append(policyCommands, policyDetails[credrequest].Command)
						// Allow only one creation command for this policy to be printed
						details.AlreadyExists = true
						policyDetails[credrequest] = details
					}
				}
			case aws.ControlPlaneCloudCredentialsRoleType:
				for i, details := range policyDetails {
					policies = append(policies, details.Name)
					if !details.AlreadyExists {
						policyCommands = append(policyCommands, details.Command)
						// Allow only one creation command for this policy to be printed
						details.AlreadyExists = true
						policyDetails[i] = details
					}
				}
			}

			// Attach policies to roles
			for _, policy := range policies {
				details, err := roles.GetPolicyDetailsByName(policyDetails, policy)
				if err != nil {
					return "", err
				}
				arn := aws.GetPolicyArn(r.Creator.Partition, r.Creator.AccountID, policy, details.Path)

				attachSharedVpcRolePolicy = awscb.NewIAMCommandBuilder().
					SetCommand(awscb.AttachRolePolicy).
					AddParam(awscb.RoleName, roleName).
					AddParam(awscb.PolicyArn, arn).
					Build()
				policyCommands = append(policyCommands, attachSharedVpcRolePolicy)
			}
		}
		commands = append(commands, createRole, attachRolePolicy)
		commands = append(commands, policyCommands...)

	}
	return awscb.JoinCommands(commands), nil
}
